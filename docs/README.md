# Budget App Documentation

This directory contains comprehensive documentation for the Budget App project.

## 📚 Documentation Index

### Architecture & Specifications

- **[Backend Specification](spec.md)** - Complete backend architecture, API design, and implementation details
- **[Architecture Diagram](img/architecture.svg)** - Visual representation of system components and data flow

### Development & Operations

- **[API Documentation](../internal/docs/)** - Auto-generated Swagger/OpenAPI documentation
- **[Database Schema](../migrations/)** - SQL migration files and schema evolution
- **[Docker Setup](../docker/README.md)** - Containerized deployment instructions

## 🏗️ System Architecture

The Budget App follows a modern microservices-inspired architecture:

```
┌─────────────────┐    HTTP/JSON    ┌─────────────────┐    SQLite    ┌─────────────────┐
│   Vue.js SPA    │ ──────────────▶ │   Go API        │ ───────────▶ │   Database      │
│   (Frontend)    │                 │   (Backend)     │              │                 │
└─────────────────┘                 └─────────────────┘              └─────────────────┘
                                              │
                                              │ Systemd Timer
                                              ▼
                                     ┌─────────────────┐
                                     │   Scheduler     │
                                     │   (Go Library)  │
                                     └─────────────────┘
```

### Key Components

1. **Frontend (Vue.js SPA)**

   - Single Page Application
   - Responsive design for mobile and desktop
   - State management with Pinia
   - API integration via Axios

2. **Backend (Go API)**

   - RESTful API built with Gin framework
   - SQLite database with sqlc code generation
   - Comprehensive input validation
   - JWT-style API key authentication

3. **Scheduler**

   - Go library function for recurring transactions
   - Systemd timer integration for reliability
   - Automatic data cleanup and backups

4. **Database**
   - SQLite for simplicity and portability
   - Optimized for single-user workloads
   - Soft delete support for data retention

## 🚀 Quick Start

### Backend Development

```bash
# Install dependencies
go mod download

# Run database migrations
make migrate

# Generate API code
make generate

# Start development server
make run
```

### Frontend Development

```bash
# Navigate to frontend directory (when created)
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

### Documentation Generation

```bash
# Generate API documentation
make api-docs

# Generate architecture diagram (requires mmdc)
make appdocs
```

## 📖 Documentation Standards

### Writing Guidelines

- Use clear, concise language
- Include code examples where relevant
- Maintain consistent formatting
- Update version numbers and dates
- Link between related documents

### File Organization

- Keep specifications in markdown format
- Use descriptive filenames
- Group related documents in subdirectories
- Maintain a clear hierarchy

### Version Control

- Commit documentation changes with code changes
- Use descriptive commit messages
- Review documentation updates in pull requests
- Tag releases with updated documentation

## 🔧 Maintenance

### Regular Updates

- Update API documentation after endpoint changes
- Refresh architecture diagrams for structural changes
- Review and update deployment instructions
- Validate code examples against current implementation

### Quality Assurance

- Test all code examples
- Verify links and references
- Check for outdated information
- Ensure consistency across documents

## 📞 Contributing

When contributing to the documentation:

1. **Follow the existing style** and format
2. **Test all examples** before committing
3. **Update related documents** when making changes
4. **Add diagrams** for complex concepts
5. **Include version information** for new features

## 📋 Documentation Checklist

Before marking documentation as complete:

- [ ] All code examples tested and working
- [ ] Links verified and functional
- [ ] Diagrams up-to-date and accurate
- [ ] Version numbers and dates current
- [ ] Grammar and spelling reviewed
- [ ] Technical accuracy verified
- [ ] Related documents updated

---

_Last updated: January 2025_
