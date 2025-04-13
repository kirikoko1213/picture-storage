import type { RequestResult } from "./request"
import request from "./request"

export interface ImageItem {
  id: number
  name: string
  url: string
  thumbnailUrl: string
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

export function apiGetImageList(directory: string, tags: string[], page: number, pageSize: number) {
  return request.post<RequestResult<ImageListResponse>>("/api/images", {
    directory,
    page,
    page_size: pageSize,
    tags,
  })
}

export function apiDeleteImages(ids: number[]) {
  return request.post<RequestResult<any>>("/api/images/delete", {
    ids,
  })
}

export function apiGetTags() {
  return request.get<RequestResult<string[]>>("/api/tags")
}

export function apiAddTags(imageIds: number[], tags: string[]) {
  return request.post<RequestResult<any>>("/api/tags/add", {
    image_ids: imageIds,
    tags,
  })
}
