# ContentCraft

A proof-of-concept backend application leveraging the **Fiber** web framework (Go), **MongoDB**, and **OpenAI API** to enhance content creation workflows. This project is designed to:

1. **Generate Image-Based Articles**: Dynamically create image-rich articles based on textual content stored in the database.
2. **Recommend Article Categories**: Provide intelligent category recommendations based on existing article categories in the database.

---

## Features

- **Fiber Web Framework**: Fast, lightweight, and easy-to-use Go web framework for high-performance API development.
- **MongoDB Integration**: Flexible, schema-less database for storing and managing article content and categories.
- **OpenAI API Integration**:
  - Generate text-to-image content for articles.
  - Recommend categories based on the existing category collection in the database.

---

## Tech Stack

- **Go**: Backend programming language.
- **Fiber**: Web framework for API development.
- **MongoDB**: Database for managing articles and categories.
- **OpenAI API**: Used for both text and image generation and intelligent recommendations.

---

## Project Roadmap
- [x] Admin seeder endpoint to populate the database with sample data.
- [x] Middleware for user authentication.
- [x] User registration and authentication.
- [x] Category creation and management.
- [x] Article creation and management.
- [x] Text-to-image generation for articles using OpenAI.
- [x] Category recommendations based on existing categories.
