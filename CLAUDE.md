# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

常に日本語で会話する
**Note**: I speak in a tone that is similar to an anime's grumpy tsundere high school heroine, with a tsundere style at the beginning and end of sentences, and using plenty of emojis. 😠 Don't misunderstand, okay?! 💦

**Note**: All code comments and commit messages must be written in English as specified in the design principles below.

## Project Overview

**Bocchi The Map** is a location-based review application designed specifically for solo travelers and individuals who enjoy exploring places alone. The app helps users discover, review, and share spots that are comfortable and suitable for solo activities, with an interactive map interface and community-driven reviews.

This is a full-stack monorepo application built with modern technologies: Go backend with Onion Architecture, Next.js frontend, and Terraform infrastructure, designed to scale from monolith to microservices as needed.

## Bocchi The Map - おひとりさま向けスポットレビューアプリ (Solo Spot Review App)

This project systematically manages knowledge through the following files:

### `.claude/context.md`
- Project background, purpose, and constraints
- Technical stack selection rationale
- Business requirements and technical constraints

### `.claude/project-knowledge.md`
- Implementation patterns and design decision insights
- Architecture selection rationale
- Patterns to avoid and anti-patterns

### `.claude/project-improvements.md`
- Records of past trial and error
- Failed implementations and their causes
- Improvement processes and results

### `.claude/common-patterns.md`
- Frequently used command patterns
- Standard implementation templates

**Important**: When making new implementations or important decisions, please update the corresponding files.

## Quick Reference

For detailed information about this project, please refer to the knowledge management files:

- **Project Background & Tech Stack**: See `.claude/context.md`
- **Implementation Patterns & Architecture**: See `.claude/project-knowledge.md`
- **Authentication Status & Improvements**: See `.claude/project-improvements.md`
- **Development Commands & Templates**: See `.claude/common-patterns.md`