import { makeAutoObservable } from "mobx";

export interface DownloadModel {
  id: string
  name: string
  percentage: number
}

class DownloadsStore {
  downloads: DownloadModel[] = [];

  constructor() {
    makeAutoObservable(this);
  }

  addDownload = (d: DownloadModel) => {
    this.downloads.push(d);
  };

  updDownload = (id: string, updatedItem: DownloadModel) => {
    const itemIndex = this.downloads.findIndex((item) => item.id === id);
    if (itemIndex !== -1) {
      this.downloads[itemIndex] = updatedItem
    } else {
      this.downloads.push(updatedItem)
    }
    let newItems = [...this.downloads]
    this.downloads = newItems
    //! rewrite?
  }

  removeDownload = (id: string) => {
    this.downloads = this.downloads.filter(d => d.id !== id)
  }
}

export const downloadStore = new DownloadsStore()
