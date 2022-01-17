export class Package {
    id : bigint
    name : string
    displayName : string
    description : string
    sourceLink : string
    creatorId : bigint
}

export function apiURL(endpoint : string) : string {
    return new URL(endpoint, import.meta.env.VITE_API_BASE_URL as string).toString()
}

export const listPackages = async (page : number, count : number) : Promise<Package[]> => {
    const res = await fetch(apiURL(`package?page=${page}&count=${count}`))
    if (!res.ok) throw new Error('Bad response')
    const items = await res.json()
    return items as Package[]
}

export const getPackage = async (packageId : bigint) : Promise<Package> => {
    const res = await fetch(apiURL(`package/${packageId}`))
    if (!res.ok) throw new Error('Bad response')
    const pack = await res.json()
    return pack as Package
}