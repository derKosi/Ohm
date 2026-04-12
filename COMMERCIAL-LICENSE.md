# Ohm Licensing

Ohm is dual-licensed to balance open-source accessibility with sustainable development.

## For Most Users: AGPL-3.0

Ohm is licensed under the [GNU Affero General Public License v3.0](LICENSE).

**You can freely:**
- Run `ohm scan` on your own machines (personal, company, enterprise)
- Use Ohm's output in your compliance reports, audits, and workflows
- Modify Ohm for your own use
- Contribute improvements back to the project

**You must:**
- Keep the license and copyright notices intact
- If you distribute modified versions, make your source code available under AGPL-3.0
- If you offer a modified version as a network service, make your source code available to those users

**You cannot:**
- Embed Ohm in a proprietary product without a commercial license
- Offer Ohm (or a modified version) as a paid service without a commercial license
- Remove or alter the license and copyright notices

## When You Need a Commercial License

You need a commercial license if you want to:

| Use Case | AGPL OK? | Commercial License |
|----------|----------|--------------------|
| Run `ohm scan` on your own machines | ✅ Yes | Not needed |
| Use Ohm output in internal reports/audits | ✅ Yes | Not needed |
| Build internal tools that call `ohm scan --json` | ✅ Yes | Not needed |
| Embed Ohm's scanner in your commercial product | ❌ No | Required |
| Offer Ohm as part of your SaaS platform | ❌ No | Required |
| Distribute modified Ohm without open-sourcing your changes | ❌ No | Required |

## Commercial License Benefits

- No AGPL copyleft obligations
- Embed Ohm in proprietary products
- Priority support and SLA
- Custom integration assistance
- Indemnification

**Contact:** mathias@kosinski.dev

## Why This Model?

Ohm is free to use for the thing it does best: scanning machines for AI software. That's valuable on its own and we want everyone to have access to it.

The commercial license funds ongoing development — new signatures, platform support, and the compliance features enterprises need. It also ensures that companies building products on top of Ohm contribute back to its development.

## FAQ

**Can I use Ohm at my company?**
Yes. Run it on every employee laptop. That's what it's for. No license needed.

**Can I build a tool that runs `ohm scan` and reads the output?**
Yes, as long as your tool is a separate program that calls Ohm (e.g., via subprocess or `--json` output). If you're importing Ohm's Go packages as libraries into your code, that's a derivative work and needs a commercial license.

**I'm a compliance consultant. Can I use Ohm at client sites?**
Yes. Running `ohm scan` at a client site is fine. If you want to embed Ohm in your consulting platform or tool, contact us for a commercial license.

**Can I contribute to Ohm?**
Yes! By contributing, you agree that your contributions can be licensed under both AGPL-3.0 and the commercial license. This is standard for dual-licensed projects and allows us to offer commercial licenses that include community contributions.
