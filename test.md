# Edit History Model (Simplified & Complete)

We model edit history using two primitives:

- **Timeline (`H`)**: an ordered list of history items `[0…N-1]` (oldest → newest).
  - `H[0]` is always a **sentinel placeholder** representing the initial design / upload.
  - All other items are **completed generations** with `{input, output, prompt}`.
- **Pointer (`p`)**: either an integer index `0…N-1` or **⊥** (null) representing the live cursor (N+1).

Together, `(H, p)` fully describe the history state.

---

## Navigation Rules

- **Rewind**
  - If `p = ⊥` → set `p = max(0, N-2)` (skip duplicate latest entry).
  - Else if `p > 0` → decrement (`p = p - 1`).
  - Else (`p = 0`) → no-op (can’t rewind past placeholder).

- **Fast-forward**
  - If `p = ⊥` → no-op (already at cursor).
  - Else if `p = N-1` → set `p = ⊥` (snap to cursor).
  - Else → increment (`p = p + 1`).

- **Counter display**
  - If `p = ⊥` → show `N / N`.
  - Else → show `(p+1) / N`.

---

## Rendering Policy

- **Image**
  - `p = ⊥` → preview latest output `H[N-1].output`.
  - `p = 0` (placeholder) → use `H[0].input`.
  - Else → use `H[p].output`.

- **Prompt**
  - `p = ⊥` or `p = N-1` → `"Type to make edits"`.
  - `p = 0` (placeholder) → `"Initial design"` or stored label.
  - Else → truncated `H[p].prompt`.

---

## Data Shape

```ts
type HistoryItem =
  | { kind: 'placeholder'; input: Asset; label?: string }
  | { kind: 'completed'; input: Asset; output: Asset; prompt: string; meta?: any };

type HistoryState = {
  H: HistoryItem[]; // H[0] is always the placeholder
  p: number | null; // null = ⊥ = cursor
};
```

---

### Transition Functions

```ts
function rewind(s: HistoryState): HistoryState {
  const N = s.H.length;
  if (s.p === null) return { ...s, p: Math.max(0, N - 2) };
  if (s.p > 0) return { ...s, p: s.p - 1 };
  return s;
}

function forward(s: HistoryState): HistoryState {
  const N = s.H.length;
  if (s.p === null) return s;
  if (s.p === N - 1) return { ...s, p: null };
  return { ...s, p: s.p + 1 };
}

function counter(s: HistoryState): string {
  const N = s.H.length;
  return s.p === null ? `${N}/${N}` : `${s.p + 1}/${N}`;
}
```

---
Edge Cases
•	Fresh session: H = [placeholder], p = ⊥.
•	Rewind guard: cannot go past p = 0.
•	Fast-forward guard: last edit skips directly to p = ⊥.
•	In-progress generation: don’t append until complete; cursor stays ⊥.
•	Failed/canceled generation: don’t append; show error instead.
•	Reloads: server always returns H[0] (placeholder); client can synthesize if missing.
•	Export/Cart/Share: skip placeholder rows when iterating H.

Optional Enhancements
•	Store a displayAssetId (e.g., before/after composite) per completed item for richer thumbnails.
•	Replace implicit ⊥ with an explicit "cursor" item if you want uniform list navigation.
•	Use parent_edit_id for branching histories (fork from any earlier item).