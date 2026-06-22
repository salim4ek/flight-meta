import { useState, type FormEvent } from 'react'
import type { SearchParams, StopsMode } from '../api'
import { PlaceInput, type PlaceValue } from './PlaceInput'

const STOPS: { value: StopsMode; label: string }[] = [
  { value: '', label: 'Любые' },
  { value: 'direct', label: 'Прямой' },
  { value: 'one', label: '1 пересадка' },
  { value: 'one_plus', label: '1+' },
]

export function SearchForm({
  onSearch,
  loading,
}: {
  onSearch: (p: SearchParams) => void
  loading: boolean
}) {
  const [origin, setOrigin] = useState<PlaceValue>({ iata: 'MOW', label: 'Москва' })
  const [destination, setDestination] = useState<PlaceValue>({ iata: 'BKK', label: 'Бангкок' })
  const [depart, setDepart] = useState('2026-04-20')
  const [ret, setRet] = useState('')
  const [passengers, setPassengers] = useState(1)
  const [passport, setPassport] = useState('RU')
  const [stops, setStops] = useState<StopsMode>('')
  const [airlines, setAirlines] = useState('')
  const [excludeAirlines, setExcludeAirlines] = useState('')
  const [selfTransfer, setSelfTransfer] = useState(true)
  const [visaFreeTransit, setVisaFreeTransit] = useState(false)
  const [hideInfeasible, setHideInfeasible] = useState(false)
  const [formError, setFormError] = useState('')

  function swap() {
    setOrigin(destination)
    setDestination(origin)
  }

  function submit(e: FormEvent) {
    e.preventDefault()
    if (!origin.iata || !destination.iata) {
      setFormError('Выберите город или страну из списка')
      return
    }
    if (origin.iata === destination.iata) {
      setFormError('Города вылета и прибытия совпадают')
      return
    }
    setFormError('')
    onSearch({
      origin: origin.iata,
      destination: destination.iata,
      depart,
      ret: ret || undefined,
      passengers,
      passport: passport.trim().toUpperCase() || undefined,
      stops,
      airlines: airlines.trim() || undefined,
      excludeAirlines: excludeAirlines.trim() || undefined,
      selfTransfer,
      visaFreeTransit,
      hideInfeasible,
    })
  }

  return (
    <form className="form" onSubmit={submit}>
      <div className="form__row">
        <div className="route-fields">
          <PlaceInput label="Откуда" placeholder="Город или страна" value={origin} onChange={setOrigin} />
          <button type="button" className="swap" onClick={swap} aria-label="Поменять местами" title="Поменять местами">
            ⇄
          </button>
          <PlaceInput label="Куда" placeholder="Город или страна" value={destination} onChange={setDestination} />
        </div>
        <label className="field">
          <span>Туда</span>
          <input type="date" value={depart} onChange={(e) => setDepart(e.target.value)} required />
        </label>
        <label className="field">
          <span>Обратно</span>
          <input type="date" value={ret} onChange={(e) => setRet(e.target.value)} />
        </label>
        <label className="field field--sm">
          <span>Пасс.</span>
          <input type="number" min={1} max={9} value={passengers} onChange={(e) => setPassengers(Number(e.target.value))} />
        </label>
        <label className="field field--sm">
          <span>Паспорт</span>
          <input value={passport} onChange={(e) => setPassport(e.target.value)} maxLength={2} placeholder="RU" />
        </label>
      </div>

      {formError && <p className="form__error">{formError}</p>}

      <div className="form__row form__row--filters">
        <div className="seg" role="group" aria-label="Пересадки">
          {STOPS.map((s) => (
            <button
              key={s.value || 'any'}
              type="button"
              className={`seg__btn${stops === s.value ? ' seg__btn--on' : ''}`}
              onClick={() => setStops(s.value)}
            >
              {s.label}
            </button>
          ))}
        </div>
        <label className="field">
          <span>Только а/к</span>
          <input value={airlines} onChange={(e) => setAirlines(e.target.value)} placeholder="TK, EK" />
        </label>
        <label className="field">
          <span>Исключить а/к</span>
          <input value={excludeAirlines} onChange={(e) => setExcludeAirlines(e.target.value)} placeholder="SU" />
        </label>
        <label className="check">
          <input type="checkbox" checked={selfTransfer} onChange={(e) => setSelfTransfer(e.target.checked)} />
          <span>Само-стыковка</span>
        </label>
        <label className="check">
          <input type="checkbox" checked={visaFreeTransit} onChange={(e) => setVisaFreeTransit(e.target.checked)} />
          <span>Только безвиз-транзит</span>
        </label>
        <label className="check">
          <input type="checkbox" checked={hideInfeasible} onChange={(e) => setHideInfeasible(e.target.checked)} />
          <span>Без невозможных стыковок</span>
        </label>
        <button className="btn btn--go" type="submit" disabled={loading}>
          {loading ? 'Ищем…' : 'Найти'}
        </button>
      </div>
    </form>
  )
}
