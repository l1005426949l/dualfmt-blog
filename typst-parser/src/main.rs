use tonic::{transport::Server, Request, Response, Status};
use parser::typst_parser_server::{TypstParser, TypstParserServer};
use parser::{ParseRequest, ParseResponse};

pub mod parser {
    tonic::include_proto!("parser");
}

// The core parsing logic from the previous step.
// I'm including it here directly for simplicity.
mod core_parser {
    use std::path::Path;
    use typst::{
        diag::{FileError, FileResult},
        eval::Tracer,
        foundations::{Bytes, Datetime},
        syntax::{FileId, Source, VirtualPath},
        text::{Font, FontBook},
        Library, World,
    };

    struct MinimalWorld {
        library: Library,
        fontbook: FontBook,
        fonts: Vec<Font>,
        source: Source,
        main: FileId,
    }

    impl MinimalWorld {
        fn new(source: &str) -> Self {
            let mut fontbook = FontBook::new();
            let mut fonts = Vec::new();

            // Use Font::iter to parse all fonts in the given data.
            for font_data in typst_assets::fonts() {
                let buffer = Bytes::from_static(font_data);
                for font in Font::iter(buffer) {
                    fontbook.push(font.info().clone());
                    fonts.push(font);
                }
            }

            let main_path = VirtualPath::new("main.typ");
            let main_id = FileId::new(None, main_path);
            let source = Source::new(main_id, source.to_string());

            Self {
                library: Library::builder().build(),
                fontbook,
                fonts,
                source,
                main: main_id,
            }
        }
    }

    impl World for MinimalWorld {
        fn library(&self) -> &Library { &self.library }
        fn main(&self) -> Source { self.source.clone() }
        fn source(&self, id: FileId) -> FileResult<Source> {
            if id == self.main { Ok(self.source.clone()) }
            else { Err(FileError::NotFound(id.vpath().as_rooted_path().to_path_buf())) }
        }
        fn book(&self) -> &FontBook { &self.fontbook }
        fn font(&self, index: usize) -> Option<Font> { self.fonts.get(index).cloned() }
        fn file(&self, id: FileId) -> FileResult<Bytes> { Err(FileError::NotFound(id.vpath().as_rooted_path().to_path_buf())) }
        fn today(&self, _offset: Option<i64>) -> Option<Datetime> { Some(Datetime::from_ymd(2024, 1, 1).unwrap()) }
    }

    pub fn parse_typst_to_html(content: &str) -> Result<String, String> {
        // NOTE: Direct HTML export is not supported in this version of the `typst`
        // library. This function will successfully compile a document to check for
        // syntax errors, but it cannot produce an HTML artifact.
        // The `typst-html` crate did not exist in a compatible way for typst v0.11.0.
        let world = MinimalWorld::new(content);
        let mut tracer = Tracer::new();
        match typst::compile(&world, &mut tracer) {
            Ok(_) => {
                // If compilation succeeds, we return an error because HTML export
                // is the part that is not supported.
                Err("Typst to HTML conversion is not supported in this version.".to_string())
            }
            Err(errors) => {
                // If compilation fails, we return the syntax errors.
                let error_str = errors.iter().map(|e| format!("{:?}", e)).collect::<Vec<_>>().join("\n");
                Err(format!("Typst compilation failed:\n{}", error_str))
            }
        }
    }
}


#[derive(Default)]
pub struct MyTypstParser;

#[tonic::async_trait]
impl TypstParser for MyTypstParser {
    async fn parse(&self, request: Request<ParseRequest>) -> Result<Response<ParseResponse>, Status> {
        let content = request.into_inner().content;

        // The result of the parse function is now always an error, either a
        // compilation failure or the "not supported" message.
        let error_message = core_parser::parse_typst_to_html(&content).unwrap_err();

        let reply = ParseResponse {
            result: Some(parser::parse_response::Result::Error(error_message)),
        };
        Ok(Response::new(reply))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let parser = MyTypstParser::default();

    println!("TypstParser server listening on {}", addr);

    Server::builder()
        .add_service(TypstParserServer::new(parser))
        .serve(addr)
        .await?;

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::core_parser::parse_typst_to_html;

    #[test]
    fn test_unsupported_feature() {
        // Test that even valid Typst content returns an error because
        // HTML export is not supported.
        let typst_content = "= Hello, Test!";
        let result = parse_typst_to_html(typst_content);
        assert!(result.is_err());
        let err_msg = result.unwrap_err();
        assert!(err_msg.contains("not supported"));
    }

    #[test]
    fn test_invalid_typst_syntax() {
        // Test that invalid syntax still produces a compilation error.
        let typst_content = "= Invalid syntax #{}";
        let result = parse_typst_to_html(typst_content);
        assert!(result.is_err());
        let err_msg = result.unwrap_err();
        assert!(err_msg.contains("compilation failed"));
    }
}
