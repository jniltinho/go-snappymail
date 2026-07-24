# Spec delta: zimbra-skin (from add-zimbra-skin)

## ADDED Requirements

### Requirement: Zimbra type-scale override
The zimbra skin SHALL override the Tailwind type-scale tokens in its scope so `text-sm`-class chrome resolves to 12px (Zimbra Classic base) without affecting other skins.

#### Scenario: Other skins unchanged
- **WHEN** the outlook or carbonio skin is active
- **THEN** `text-sm` still resolves to 14px

### Requirement: Skin-styled dropdown menus
Dropdown menus SHALL render as Zimbra Classic menus under the zimbra skin: white panel, 1px `#bfbfbf` border, subtle shadow, 12px items, `#CCE5F3` hover, squared.

#### Scenario: Menu chrome
- **WHEN** any dropdown opens under the zimbra skin
- **THEN** its computed panel/border/hover values match the reference menu
