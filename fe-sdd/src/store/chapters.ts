import { makeAutoObservable } from 'mobx'
import { Volume, Chapter } from '../api/models'

class ChapterStore {
  volToChapters: Map<string, Set<string>> = new Map()

  constructor() {
    makeAutoObservable(this)
  }

  init(volumes: Volume[]) {
    volumes.forEach((v) => {
      const chSet: Set<string> = new Set()
      v.Chapters.forEach((c) => {
        chSet.add(c.UID)
      })
      this.volToChapters.set(v.UID, chSet)
    })
  }

  uncheckVolume(vID: string) {
    this.volToChapters = this.volToChapters.set(vID, new Set())
  }

  checkVolume(v: Volume) {
    const chSet: Set<string> = new Set()
    v.Chapters.forEach((c) => {
      chSet.add(c.UID)
    })
    this.volToChapters = this.volToChapters.set(v.UID, chSet)
  }

  addChapter(vID: string, cID: string) {
    const chapters = this.volToChapters.get(vID)
    if (!chapters) {
      console.log('this should not happened')
      return
    }
    this.volToChapters = this.volToChapters.set(vID, chapters.add(cID))
  }

  removeChapter(vID: string, cID: string): boolean {
    const chapters = this.volToChapters.get(vID)
    if (!chapters) {
      return false
    }
    if (!chapters.delete(cID)) {
      return false
    }
    this.volToChapters = this.volToChapters.set(vID, chapters)
    console.log(this.volToChapters.get(vID))
    return true
  }
}

export const chapterStore = new ChapterStore()
