# Contributing Guidelines

Thank you for considering contributing to this project! This guide outlines the contribution process and expectations.

## How to Contribute

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/your-username/devin-terraform.git
cd devin-terraform

# Install dependencies
go mod tidy

# Build
make build

# Test
make test
```

## Pull Request Guidelines

1. Code should be formatted according to Go standard style
2. All tests should pass
3. New features should include corresponding tests
4. Include a clear PR description explaining the changes

## Commit Message Guidelines

Please follow this format for commit messages:

```
feat: Add new feature
fix: Bug fix
docs: Documentation-only changes
style: Changes that do not affect code behavior (formatting, etc.)
refactor: Code changes that neither fix bugs nor add features
perf: Changes that improve performance
test: Adding or correcting tests
chore: Changes to build process or tools
```

## Code Review Process

Pull requests are reviewed by maintainers. To facilitate the review process:

1. Keep pull requests small and focused on a clear objective
2. Add comments to your code to explain complex sections
3. Respond to feedback from reviews

## License

Contributions to the project are licensed under the same license as the project (MIT).
