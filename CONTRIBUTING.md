# Contributing to Ohm

Thank you for your interest in contributing to Ohm!

## Developer Certificate of Origin (DCO)

By contributing to this project, you agree to the [Developer Certificate of Origin](https://developercertificate.org/):

> Developer Certificate of Origin
> Version 1.1
>
> Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
>
> Everyone is permitted to copy and distribute verbatim copies of this
> license document, but changing it is not allowed.
>
> Developer's Certificate of Origin 1.1
>
> By making a contribution to this project, I certify that:
>
> (a) The contribution was created in whole or in part by me and I
>     have the right to submit it under the open source license
>     indicated in the file; or
>
> (b) The contribution is based upon previous work that, to the best
>     of my knowledge, is covered under an appropriate open source
>     license and I have the right under that license to submit that
>     work with modifications, whether created in whole or in part
>     by me, under the same open source license (unless I am
>     permitted to submit under a different license), as indicated
>     in the file; or
>
> (c) The contribution was provided directly to me by some other
>     person who certified (a), (b) or (c) and I have not modified
>     it.
>
> (d) I understand and agree that this project and the contribution
>     are public and that a record of the contribution (including all
>     personal information I submit with it, including my sign-off) is
>     maintained indefinitely and may be redistributed consistent with
>     this project or the open source license(s) involved.

### How to sign off

Every commit must include a `Signed-off-by` line:

```bash
git commit -s -m "your commit message"
```

This adds `Signed-off-by: Your Name <your@email.com>` to the commit message.

PRs without sign-off on every commit will not be merged.

## License

By contributing, you agree that your contributions will be licensed under both:

- **AGPL-3.0-or-later** (the project's open-source license)
- **The commercial license** (used for dual-licensing the project)

This is standard for dual-licensed open-source projects and allows us to offer commercial licenses that include community contributions.

## Adding New Signatures

When adding a new AI tool signature:

1. Verify the tool is real (web search — no made-up entries)
2. Add the detection entry in `internal/scanner/<category>.go`
3. Include all 3 platform uninstall commands (linux, macos, windows)
4. Choose the appropriate risk level:
   - `RiskSafe` — no sensitive data
   - `RiskCaution` — config files that may contain preferences
   - `RiskDanger` — directories likely to contain API keys or credentials
5. Add the tool to `docs/SIGNATURES.md`

## Code Style

- `gofmt` formatted
- `go vet` clean
- No external dependencies without discussion (we want to stay lean)

## Questions?

Open an issue or start a discussion on GitHub.
