export interface MangaDex {
  Info: MangaInfo
  Volumes: Volume[]
}

export interface MangaInfo {
  Title: string
  Authors: string[]
  Artists: string[]
  ID: string
}

export interface Volume {
  UID: string
  Info: VolumeInfo
  Chapters: Chapter[]
  Cover: any
  CoverPath: string
}

export interface VolumeInfo {
  Identifier: string
}

export interface Chapter {
  UID: string
  Info: ChapterInfo
  Pages: Pages
}

export interface ChapterInfo {
  Title: string
  Views: number
  Language: string
  GroupNames: string[]
  Published: string
  ID: string
  Identifier: string
  VolumeIdentifier: string
}

export interface Pages {}

export interface StoredManga {
  Title: string
  Cover: string
  VolumeCount: number
  IsUpscaled: boolean
  AverageMegapixels: number
}
