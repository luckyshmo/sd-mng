import axios from "axios";

const address = process.env.REACT_APP_API_URL!;

export interface LoraInfo {
  name: string
  token: string
}

export async function getLoraInfos(): Promise<LoraInfo[]> {
  try {
    const response = await axios.get(address+"info/lora");
    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }
    return response.data
  } catch (error: any) {
    //! how to handle this properly?
    console.error(`Error on getting loras: ${error.message}`);
    return []
  }
}

export async function downloadFile(url: string, folder: string): Promise<string> {
  const urlParams = new URLSearchParams({
    url: url,
    folder: folder,
  });

  try {
    const response = await axios.post(address + "?" + urlParams.toString());

    if (response.status !== 200) {
      throw new Error(`HTTP error! Status: ${response.status}`);
    }

    console.log(response.data);
    return ""
  } catch (error: any) {
    console.error(`Error: ${error.message}`);
    return error.message
  }
}