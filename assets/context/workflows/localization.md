---
ctx:
  tags: [workflow, i18n, l10n]
  purpose: Record how this project handles translation, locale-sensitive formatting, and string management.
  use_when:
    - adding user-facing strings
    - changing translation workflows or locale behavior
  fill:
    - required: supported locales and source locale
    - required: string storage and translation workflow
    - optional: testing and formatting rules
---

# Localization

> **Fill this in.** This file is project-specific and cannot be generalized.
>
> Document the localization (l10n) and internationalization (i18n) workflow for this project. This helps AI assistants add user-facing strings correctly and avoid bypassing the localization system.
>
> Suggested structure:
> - **Supported locales:** which languages/regions are supported and which is the source locale
> - **Tooling:** what i18n library or framework is used (e.g. i18next, gettext, ICU)
> - **String management:** how strings are added, where they live, and how translations are managed
> - **Translation workflow:** how new strings get translated (internal team, vendor, community)
> - **Locale-sensitive formats:** how dates, numbers, currencies, and plurals are handled
> - **Testing:** how to test a specific locale locally
