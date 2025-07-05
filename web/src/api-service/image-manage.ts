import type { RequestResult } from "./request"
import request from "./request"

export interface ImageItem {
  id: number
  name: string
  url: string
  thumbnailUrl: string
  tags: string[]
  size: number
  created_at: string
}

export interface ImageListResponse {
  list: ImageItem[]
  total: number
}

export interface TagItem {
  name: string
  count: number
}

export interface TagListResponse {
  list: TagItem[]
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
  return request.delete<RequestResult<any>>("/api/images", {
    ids,
  })
}

export function apiGetTags() {
  return request.get<RequestResult<string[]>>("/api/tags")
}

export function apiGetTagDetails() {
  return request.get<RequestResult<TagListResponse>>("/api/tags/details")
}

export function apiCreateTag(name: string) {
  return request.post<RequestResult<any>>("/api/tags", {
    name,
  })
}

export function apiDeleteTag(name: string) {
  return request.delete<RequestResult<any>>("/api/tags", {
    name,
  })
}

export function apiUpdateTag(oldName: string, newName: string) {
  return request.put<RequestResult<any>>("/api/tags", {
    old_name: oldName,
    new_name: newName,
  })
}

export function apiAddTags(imageIds: number[], tags: string[]) {
  return request.post<RequestResult<any>>("/api/images/tags", {
    image_ids: imageIds,
    tags,
  })
}
