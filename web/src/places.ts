// Human-friendly place picker data. Users search by city or country name; the
// backend still works in IATA codes, so each entry maps to one. Picking a
// country selects its main city automatically.

export interface City {
  name: string
  iata: string
  country: string
  popular?: boolean
}

export interface Country {
  name: string
  city: string // main city name
  iata: string // main city IATA
}

const CITIES: City[] = [
  { name: 'Москва', iata: 'MOW', country: 'Россия', popular: true },
  { name: 'Санкт-Петербург', iata: 'LED', country: 'Россия', popular: true },
  { name: 'Сочи', iata: 'AER', country: 'Россия', popular: true },
  { name: 'Екатеринбург', iata: 'SVX', country: 'Россия' },
  { name: 'Новосибирск', iata: 'OVB', country: 'Россия' },
  { name: 'Казань', iata: 'KZN', country: 'Россия' },
  { name: 'Краснодар', iata: 'KRR', country: 'Россия' },
  { name: 'Красноярск', iata: 'KJA', country: 'Россия' },
  { name: 'Уфа', iata: 'UFA', country: 'Россия' },
  { name: 'Калининград', iata: 'KGD', country: 'Россия' },
  { name: 'Владивосток', iata: 'VVO', country: 'Россия' },
  { name: 'Минеральные Воды', iata: 'MRV', country: 'Россия' },
  { name: 'Махачкала', iata: 'MCX', country: 'Россия' },
  { name: 'Бангкок', iata: 'BKK', country: 'Таиланд', popular: true },
  { name: 'Пхукет', iata: 'HKT', country: 'Таиланд' },
  { name: 'Стамбул', iata: 'IST', country: 'Турция', popular: true },
  { name: 'Анталья', iata: 'AYT', country: 'Турция' },
  { name: 'Дубай', iata: 'DXB', country: 'ОАЭ', popular: true },
  { name: 'Абу-Даби', iata: 'AUH', country: 'ОАЭ' },
  { name: 'Доха', iata: 'DOH', country: 'Катар' },
  { name: 'Пекин', iata: 'PEK', country: 'Китай' },
  { name: 'Шанхай', iata: 'PVG', country: 'Китай' },
  { name: 'Гуанчжоу', iata: 'CAN', country: 'Китай' },
  { name: 'Гонконг', iata: 'HKG', country: 'Гонконг' },
  { name: 'Дели', iata: 'DEL', country: 'Индия' },
  { name: 'Мумбаи', iata: 'BOM', country: 'Индия' },
  { name: 'Коломбо', iata: 'CMB', country: 'Шри-Ланка' },
  { name: 'Мале', iata: 'MLE', country: 'Мальдивы' },
  { name: 'Куала-Лумпур', iata: 'KUL', country: 'Малайзия' },
  { name: 'Сингапур', iata: 'SIN', country: 'Сингапур' },
  { name: 'Белград', iata: 'BEG', country: 'Сербия' },
  { name: 'Ереван', iata: 'EVN', country: 'Армения' },
  { name: 'Тбилиси', iata: 'TBS', country: 'Грузия' },
  { name: 'Баку', iata: 'GYD', country: 'Азербайджан' },
  { name: 'Ташкент', iata: 'TAS', country: 'Узбекистан' },
  { name: 'Алматы', iata: 'ALA', country: 'Казахстан' },
  { name: 'Астана', iata: 'NQZ', country: 'Казахстан' },
  { name: 'Минск', iata: 'MSQ', country: 'Беларусь' },
  { name: 'Бишкек', iata: 'FRU', country: 'Кыргызстан' },
  { name: 'Душанбе', iata: 'DYU', country: 'Таджикистан' },
]

const COUNTRIES: Country[] = [
  { name: 'Россия', city: 'Москва', iata: 'MOW' },
  { name: 'Таиланд', city: 'Бангкок', iata: 'BKK' },
  { name: 'Турция', city: 'Стамбул', iata: 'IST' },
  { name: 'ОАЭ', city: 'Дубай', iata: 'DXB' },
  { name: 'Катар', city: 'Доха', iata: 'DOH' },
  { name: 'Китай', city: 'Пекин', iata: 'PEK' },
  { name: 'Индия', city: 'Дели', iata: 'DEL' },
  { name: 'Шри-Ланка', city: 'Коломбо', iata: 'CMB' },
  { name: 'Мальдивы', city: 'Мале', iata: 'MLE' },
  { name: 'Малайзия', city: 'Куала-Лумпур', iata: 'KUL' },
  { name: 'Сингапур', city: 'Сингапур', iata: 'SIN' },
  { name: 'Сербия', city: 'Белград', iata: 'BEG' },
  { name: 'Армения', city: 'Ереван', iata: 'EVN' },
  { name: 'Грузия', city: 'Тбилиси', iata: 'TBS' },
  { name: 'Азербайджан', city: 'Баку', iata: 'GYD' },
  { name: 'Узбекистан', city: 'Ташкент', iata: 'TAS' },
  { name: 'Казахстан', city: 'Алматы', iata: 'ALA' },
  { name: 'Беларусь', city: 'Минск', iata: 'MSQ' },
]

export interface Place {
  kind: 'city' | 'country'
  iata: string
  name: string
  sub: string
  label: string // what goes into the input
}

function cityPlace(c: City): Place {
  return { kind: 'city', iata: c.iata, name: c.name, sub: c.country, label: c.name }
}

function countryPlace(c: Country): Place {
  return { kind: 'country', iata: c.iata, name: c.name, sub: `Страна · ${c.city}`, label: `${c.name} · ${c.city}` }
}

// searchPlaces returns city + country matches. With an empty query it returns a
// few popular cities (shown when the field is focused).
export function searchPlaces(query: string, showPopular: boolean): Place[] {
  const q = query.trim().toLowerCase()
  if (!q) {
    return showPopular ? CITIES.filter((c) => c.popular).map(cityPlace) : []
  }
  const cities = CITIES.filter(
    (c) => c.name.toLowerCase().startsWith(q) || c.name.toLowerCase().includes(q) || c.iata.toLowerCase() === q,
  )
    .slice(0, 6)
    .map(cityPlace)
  const countries = COUNTRIES.filter((c) => c.name.toLowerCase().includes(q)).slice(0, 3).map(countryPlace)
  return [...cities, ...countries].slice(0, 8)
}
