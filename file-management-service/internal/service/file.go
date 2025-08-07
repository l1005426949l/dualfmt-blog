package service

import (
	"bytes"
	"context"
	"io"

	pb "file-management-service/api/file/v1"
	"file-management-service/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type FileService struct {
	pb.UnimplementedFileServer

	uc  *biz.FileUsecase
	log *log.Helper
}

func NewFileService(uc *biz.FileUsecase, logger log.Logger) *FileService {
	return &FileService{uc: uc, log: log.NewHelper(logger)}
}

func (s *FileService) UploadFile(stream pb.File_UploadFileServer) error {
	req, err := stream.Recv()
	if err != nil {
		s.log.Errorf("failed to receive first upload message: %v", err)
		return pb.ErrorBadRequest("failed to receive first upload message")
	}

	info := req.GetInfo()
	if info == nil {
		return pb.ErrorBadRequest("first message must contain file info")
	}

	s.log.Infof("receiving file upload: %s", info.FileName)

	var buf bytes.Buffer
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.log.Errorf("failed to receive upload chunk: %v", err)
			return pb.ErrorInternal("failed to receive upload chunk")
		}
		chunk := req.GetChunkData()
		if _, err := buf.Write(chunk); err != nil {
			s.log.Errorf("failed to write chunk to buffer: %v", err)
			return pb.ErrorInternal("failed to write chunk to buffer")
		}
	}

	file := &biz.File{
		Name:    info.FileName,
		Content: &buf,
		Size:    int64(buf.Len()),
	}

	filePath, err := s.uc.UploadFile(stream.Context(), file, info.FileType)
	if err != nil {
		s.log.Errorf("failed to upload file: %v", err)
		return pb.ErrorInternal("failed to upload file")
	}

	return stream.SendAndClose(&pb.UploadFileResponse{
		FilePath: filePath,
		Size:     uint64(file.Size),
	})
}

func (s *FileService) DownloadFile(req *pb.DownloadFileRequest, stream pb.File_DownloadFileServer) error {
	s.log.Infof("handling download request for: %s", req.FilePath)

	file, err := s.uc.DownloadFile(stream.Context(), req.FilePath)
	if err != nil {
		s.log.Errorf("failed to download file from biz layer: %v", err)
		return pb.ErrorNotFound("file not found")
	}
	// The reader from minio needs to be closed.
	if closer, ok := file.Content.(io.Closer); ok {
		defer closer.Close()
	}

	buf := make([]byte, 1024*4) // 4KB chunks
	for {
		n, err := file.Content.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			s.log.Errorf("failed to read file chunk: %v", err)
			return pb.ErrorInternal("failed to read file chunk")
		}
		if err := stream.Send(&pb.DownloadFileResponse{ChunkData: buf[:n]}); err != nil {
			s.log.Errorf("failed to send file chunk: %v", err)
			return pb.ErrorInternal("failed to send file chunk")
		}
	}
	return nil
}
