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

            for font_data in typst_assets::fonts() {
                let buffer = Bytes::from_static(font_data);
                let face_count = ttf_parser::fonts_in_collection(&buffer).unwrap_or(1);
                for i in 0..face_count {
                    if let Some(font) = Font::new(buffer.clone(), i) {
                        fontbook.push(font.info().clone());
                        fonts.push(font);
                    }
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
        let world = MinimalWorld::new(content);
        let mut tracer = Tracer::new();
        match typst::compile(&world, &mut tracer) {
            Ok(document) => Ok(typst_html::html(&document)),
            Err(errors) => {
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

        match core_parser::parse_typst_to_html(&content) {
            Ok(html) => {
                let reply = ParseResponse {
                    result: Some(parser::parse_response::Result::HtmlContent(html)),
                };
                Ok(Response::new(reply))
            }
            Err(e) => {
                 let reply = ParseResponse {
                    result: Some(parser::parse_response::Result::Error(e)),
                };
                Ok(Response::new(reply))
            }
        }
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
    fn test_parse_simple_heading() {
        let typst_content = "= Hello, Test!";
        let result = parse_typst_to_html(typst_content);
        assert!(result.is_ok());
        let html = result.unwrap();
        // A level 1 heading should be rendered as an h1 tag.
        // We check for `>Hello, Test!<` to be flexible with attributes on the h1 tag.
        assert!(html.contains("<h1>Hello, Test!</h1>"));
    }

    #[test]
    fn test_parse_math() {
        let typst_content = "Here is some math: $a + b = c$";
        let result = parse_typst_to_html(typst_content);
        assert!(result.is_ok());
        let html = result.unwrap();
        // Math is typically rendered with special tags or classes.
        // We'll check for the presence of the formula itself.
        assert!(html.contains("a + b = c"));
    }

    #[test]
    fn test_invalid_typst_syntax() {
        let typst_content = "= Invalid syntax #{}";
        let result = parse_typst_to_html(typst_content);
        assert!(result.is_err());
    }
}
