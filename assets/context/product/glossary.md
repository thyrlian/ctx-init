---
ctx:
  tags: [product, global]
  purpose: Define domain terms and acronyms so AI and humans interpret project language consistently.
  use_when:
    - a term appears domain-specific or ambiguous
    - writing docs, comments, prompts, or user-facing copy
  fill:
    - required: term
    - required: definition in this project's context
    - optional: synonyms, related terms, and common confusion
---

# Glossary

> Project-specific. Start with terms that are non-obvious, overloaded, or frequently confused.
>
> From a Domain-Driven Design perspective, this file is the entry point to the
> project's ubiquitous language: the shared vocabulary that product, design,
> engineering, and AI systems should use consistently.
>
> Useful guidance:
> - Prefer domain terms over implementation slang.
> - Record terms that look similar but are not interchangeable.
> - If the same word means different things in different bounded contexts,
>   document those meanings separately instead of forcing one global definition.
> - When helpful, note the code-facing name used for the same concept so readers
>   can connect domain language to implementation without collapsing the two.
> - If a historical term is still present in code or conversation but is no
>   longer preferred, mark it explicitly to prevent future ambiguity.
