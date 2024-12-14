# GeminiContentGen

A proof-of-concept backend application leveraging the **Fiber** web framework (Go), **MongoDB**, and **Gemini LLM API** to enhance content creation workflows. This project is designed to:

1. **Generate Image-Based Articles**: Dynamically create image-rich articles based on textual content stored in the database.
2. **Recommend Article Categories**: Provide intelligent category recommendations based on existing article categories in the database.

---

## Features

- **Fiber Web Framework**: Fast, lightweight, and easy-to-use Go web framework for high-performance API development.
- **MongoDB Integration**: Flexible, schema-less database for storing and managing article content and categories.
- **Gemini LLM Integration**:
  - Generate images for articles based on their content.
  - Recommend categories based on the existing category collection in the database.

---

## Tech Stack

- **Go**: Backend programming language.
- **Fiber**: Web framework for API development.
- **MongoDB**: Database for managing articles and categories.
- **Gemini LLM API**: For generating images and category recommendations.

---

## Project roadmap
- [x] Admin seeder endpoint to populate the database with sample data.
- [x] User registration and authentication.
- [x] Category creation and management.
- [ ] Article creation and management.
- [ ] Image generation for articles.
- [ ] Category recommendations based on existing categories.