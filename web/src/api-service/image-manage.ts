import type { RequestResult } from "./request"
import request from "./request"

export interface ImageItem {
  id: number
  name: string
  url: string
  size: number
  created_at: string
}

export interface ImageListResponse {
  list: ImageItem[]
  total: number
}

export function apiGetDirectoryList() {
  return request.get<RequestResult<string[]>>("/api/directory")
}

export function apiUploadImage(directory: string, file: File) {
  return request.post<RequestResult<any>>("/api/upload", {
    directory,
    file,
  })
}

export function apiGetImageList(directory: string, page: number, pageSize: number) {
  return request.get<RequestResult<ImageListResponse>>("/api/images", {
    directory,
    page,
    page_size: pageSize,
  })
}

export function apiDeleteImages(ids: number[]) {
  return request.post<RequestResult<any>>("/api/images/delete", {
    ids,
  })
}
