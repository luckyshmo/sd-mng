export interface Manga {
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
  Info: VolumeInfo
  Chapters: Chapter[]
  Cover: any
  CoverPath: string
}

export interface VolumeInfo {
  Identifier: string
}

export interface Chapter {
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
