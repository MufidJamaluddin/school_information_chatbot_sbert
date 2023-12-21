export interface IPagination<T> {
  start: number
  end: number
  length: number
  data: Array<T>
}

export interface SearchParam {
  page: number
  size: number
  keyword: string
}