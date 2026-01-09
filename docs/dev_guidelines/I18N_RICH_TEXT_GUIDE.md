# i18n Rich Text Guide

This guide explains how to handle rich text (bold, italics, links, etc.) in internationalized (i18n) translations following Vue i18n best practices.

## Table of Contents
- [Problem: HTML in Translations](#problem-html-in-translations)
- [Solution: i18n-t Component](#solution-i18n-t-component)
- [Best Practices](#best-practices)
- [Examples](#examples)
- [Anti-Patterns](#anti-patterns)

## Problem: HTML in Translations

### ❌ Don't Do This

Vue i18n **blocks HTML in translation strings** by default for security reasons (XSS prevention):

```json
// ❌ BAD: HTML in translations
{
  "message": "Are you sure you want to delete <strong>{name}</strong>?"
}
```

```vue
<!-- ❌ BAD: Using v-html with translations -->
<p v-html="t('message', { name: 'John' })" />
```

**Problems:**
- Vue i18n will throw an error: `Detected HTML in message`
- Security risk: opens door to XSS attacks
- Mixing content with presentation
- Not maintainable

## Solution: i18n-t Component

### ✅ The Right Way

Use the `<i18n-t>` component with **named slots** for rich text interpolation:

**1. Keep Translation Pure (No HTML)**
```json
{
  "message": "Are you sure you want to delete {name}?"
}
```

**2. Use i18n-t Component with Slots**
```vue
<i18n-t keypath="message" tag="span">
  <template #name>
    <strong>{{ userName }}</strong>
  </template>
</i18n-t>
```

## Best Practices

### 1. Always Remove HTML from Translations

**Before:**
```json
{
  "warning": "<strong>Warning:</strong> This action is permanent for {user}."
}
```

**After:**
```json
{
  "warning": "Warning: This action is permanent for {user}."
}
```

### 2. Use i18n-t for Dynamic Content with Styling

**Component Template:**
```vue
<AlertDescription>
  <i18n-t keypath="features.oauth_clients.dialogs.delete.description" tag="span">
    <template #name>
      <strong>{{ client.name }}</strong>
    </template>
  </i18n-t>
</AlertDescription>
```

**Translation:**
```json
{
  "features": {
    "oauth_clients": {
      "dialogs": {
        "delete": {
          "description": "Are you sure you want to delete {name}? This cannot be undone."
        }
      }
    }
  }
}
```

### 3. Wrap Static Emphasis in Components

For text that doesn't have dynamic interpolation but needs emphasis:

```vue
<p class="text-sm text-yellow-800">
  <strong>{{ t('features.oauth_clients.secretDisplay.warning') }}</strong>
</p>
```

```json
{
  "warning": "Important: Save this secret now - it won't be shown again!"
}
```

### 4. i18n-t Component Props

- **`keypath`**: Translation key path (required)
- **`tag`**: HTML tag to wrap content (default: `'span'`)
- **`plural`**: Number for pluralization (optional)
- **`scope`**: Component scope (optional)

## Examples

### Example 1: Delete Confirmation

**Component (OAuthClientDeleteDialog.vue):**
```vue
<AlertDialogDescription>
  <i18n-t keypath="features.oauth_clients.dialogs.delete.description" tag="span">
    <template #name>
      <strong>{{ client.name }}</strong>
    </template>
  </i18n-t>
</AlertDialogDescription>
```

**Translation (en-US.json):**
```json
{
  "features": {
    "oauth_clients": {
      "dialogs": {
        "delete": {
          "description": "Are you sure you want to delete {name}? This action cannot be undone. All applications using this client will stop working."
        }
      }
    }
  }
}
```

**Translation (id-ID.json):**
```json
{
  "features": {
    "oauth_clients": {
      "dialogs": {
        "delete": {
          "description": "Apakah Anda yakin ingin menghapus {name}? Tindakan ini tidak dapat dibatalkan. Semua aplikasi yang menggunakan klien ini akan berhenti bekerja."
        }
      }
    }
  }
}
```

### Example 2: Multiple Interpolations

**Component:**
```vue
<i18n-t keypath="user.greeting" tag="p">
  <template #name>
    <strong>{{ user.name }}</strong>
  </template>
  <template #role>
    <Badge>{{ user.role }}</Badge>
  </template>
</i18n-t>
```

**Translation:**
```json
{
  "user": {
    "greeting": "Welcome {name}! You are logged in as {role}."
  }
}
```

### Example 3: Links in Text

**Component:**
```vue
<i18n-t keypath="terms.accept" tag="p">
  <template #link>
    <NuxtLink to="/terms" class="underline">
      {{ t('terms.linkText') }}
    </NuxtLink>
  </template>
</i18n-t>
```

**Translation:**
```json
{
  "terms": {
    "accept": "By continuing, you agree to our {link}.",
    "linkText": "Terms of Service"
  }
}
```

### Example 4: Multiple Named Slots with Different Styles

**Component:**
```vue
<i18n-t keypath="pricing.offer" tag="div" class="text-lg">
  <template #price>
    <span class="text-2xl font-bold text-green-600">{{ price }}</span>
  </template>
  <template #discount>
    <Badge variant="destructive">{{ discount }}</Badge>
  </template>
  <template #period>
    <span class="text-sm text-muted-foreground">{{ period }}</span>
  </template>
</i18n-t>
```

**Translation:**
```json
{
  "pricing": {
    "offer": "Get it for {price} with {discount} off for the first {period}!"
  }
}
```

## Anti-Patterns

### ❌ Don't Use v-html with Translations

```vue
<!-- ❌ BAD: XSS risk, vue-i18n will error -->
<div v-html="t('message', { name: clientName })" />
```

### ❌ Don't Put HTML in Translation Files

```json
// ❌ BAD: Will cause errors
{
  "message": "Click <a href='/terms'>here</a> to continue"
}
```

### ❌ Don't Concatenate Translations

```vue
<!-- ❌ BAD: Breaks i18n, hard to translate -->
<p>{{ t('prefix') }} <strong>{{ value }}</strong> {{ t('suffix') }}</p>
```

### ❌ Don't Use String Interpolation for Rich Text

```vue
<!-- ❌ BAD: Can't be translated properly -->
<p>`${t('message')} <strong>${name}</strong>`</p>
```

## Why This Approach?

### Security
- ✅ **No XSS vulnerabilities**: Vue components are safe, v-html is not
- ✅ **Content Security Policy friendly**: No inline HTML parsing

### i18n Best Practices
- ✅ **Translator-friendly**: Translators see clean text, not HTML
- ✅ **Context-aware**: Different languages can rearrange {placeholders}
- ✅ **RTL support**: Works seamlessly with right-to-left languages

### Maintainability
- ✅ **Separation of concerns**: Content (translations) vs. presentation (Vue components)
- ✅ **Type-safe**: Vue components provide better IDE support
- ✅ **Reusable**: Components can be styled/themed independently

### Examples of Language Word Order Differences

**English:** "Are you sure you want to delete {name}?"
**Japanese:** "{name}を削除してもよろしいですか？" (Literal: "{name} delete are you sure?")
**Arabic:** "هل أنت متأكد من حذف {name}؟" (Right-to-left, different structure)

The `i18n-t` component handles these variations automatically!

## Quick Reference

| Scenario | Solution | Example |
|----------|----------|---------|
| Dynamic content needs emphasis | `<i18n-t>` with slot | `<template #name><strong>{{ name }}</strong></template>` |
| Static text needs emphasis | Wrap in component | `<strong>{{ t('key') }}</strong>` |
| Links in text | `<i18n-t>` with link slot | `<template #link><NuxtLink>...</NuxtLink></template>` |
| Multiple styled parts | `<i18n-t>` with multiple slots | See Example 4 above |
| Plain text only | Direct interpolation | `{{ t('key') }}` |

## Migration Checklist

When converting v-html translations to i18n-t:

- [ ] Remove HTML tags from translation JSON files
- [ ] Replace `v-html="t(...)"` with `<i18n-t>`
- [ ] Add named slots for each `{placeholder}`
- [ ] Apply styling/components in slot templates
- [ ] Test with all supported locales
- [ ] Verify no vue-i18n console errors
- [ ] Run ESLint (should pass with no v-html warnings)

## Resources

- [Vue I18n Official Docs - Component Interpolation](https://vue-i18n.intlify.dev/guide/advanced/component.html)
- [Nuxt I18n Module](https://i18n.nuxtjs.org/)
- [OWASP XSS Prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross_Site_Scripting_Prevention_Cheat_Sheet.html)

## Related Files

- OAuth Client implementation examples:
  - `frontend/app/components/features/oauth-client/OAuthClientDeleteDialog.vue`
  - `frontend/app/components/features/oauth-client/OAuthClientRevealDialog.vue`
  - `frontend/app/components/features/oauth-client/OAuthClientSecretDisplay.vue`
