# Learning In Public - A Go-Based Blog

##  What it is
I'm building this blog to share my learning of golang and other technologies in public. I wanted to learn more about web development with Go (I know not the most popular). This is mainly a place for me to start writing and share those writings with other people.

## Learning Goals:

### General / Why:
- Learning how to write better.
    - I think that writing is one of the most important skills someone can have. I think that the way someone writes is the best way to see how they think

### What I'm Learning From This App
- Learning more about Go and Building websites with Go
- Learning more about how to build an MVP and iterate on a project. Trying to get better at planning and creating new things from scratch and evolving them over time

## Tech Stack
- Go -- I want to get better at Go and be a Gopher
- More go -- Standard library is the way (star wars meme)

## Architecture
- Markdown files stored in /posts directory
- Go server reads files on startup
- Templ renders HTML
- Server reads markdown files and parses frontmatter for metadata

## Features (v1 Scope)
- Home page (`/`) displays about section and 3 most recent blog posts as cards
- Blog listing page (`/blog`) shows all posts chronologically
- Individual post pages (`/blog/{slug}`) render full markdown content
- Posts are stored as markdown files on the server
- New posts added via SFTP to server directory

## Future Ideas
- SQLite DB (If I want comments or something)
- Typst support (instead of markdown)
- Portfolio page with Github api integration

## Setup and Running
<!-- TODO -->
