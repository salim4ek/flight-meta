import { useEffect, useRef, useState } from 'react'
import { searchPlaces, type Place } from '../places'

export interface PlaceValue {
  iata: string
  label: string
}

export function PlaceInput({
  label,
  placeholder,
  value,
  onChange,
}: {
  label: string
  placeholder: string
  value: PlaceValue
  onChange: (v: PlaceValue) => void
}) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLLabelElement>(null)

  useEffect(() => {
    function onDoc(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', onDoc)
    return () => document.removeEventListener('mousedown', onDoc)
  }, [])

  const results = searchPlaces(value.label, open)

  function pick(p: Place) {
    onChange({ iata: p.iata, label: p.label })
    setOpen(false)
  }

  return (
    <label className="field place" ref={ref}>
      <span>{label}</span>
      <input
        value={value.label}
        placeholder={placeholder}
        autoComplete="off"
        onFocus={() => setOpen(true)}
        onChange={(e) => {
          // typing clears the selected code until a list item is picked
          onChange({ iata: '', label: e.target.value })
          setOpen(true)
        }}
      />
      {open && results.length > 0 && (
        <ul className="place__menu">
          {results.map((p) => (
            <li key={p.kind + p.iata + p.name}>
              <button
                type="button"
                className={`place__opt place__opt--${p.kind}`}
                onMouseDown={(e) => e.preventDefault()}
                onClick={() => pick(p)}
              >
                <span className="place__name">{p.name}</span>
                <span className="place__sub">{p.sub}</span>
              </button>
            </li>
          ))}
        </ul>
      )}
    </label>
  )
}
