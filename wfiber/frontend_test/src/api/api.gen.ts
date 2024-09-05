// this file was generated DO NOT EDIT
export interface SimpleAPIResponse {
	success: boolean
	error: string
}
export interface SimpleAPIResponseSearchCritCondition {
	field: 'success' | 'error'
	condition: 'equals' | 'not equals' | 'contains' | 'not contains' | 'greater than' | 'lesser than'
	value: any
}
export interface SimpleAPIResponseWithData<T> {
	success: boolean
	error: string
	data: T
}
export interface Tests {
	ID: number
	name: string
}
export interface TestsSearchCritCondition {
	field: 'ID' | 'name'
	condition: 'equals' | 'not equals' | 'contains' | 'not contains' | 'greater than' | 'lesser than'
	value: any
}
import http from './http'
export const GetTests = (query?: string, filter?: TestsSearchCritCondition[]): Promise<SimpleAPIResponseWithData<Tests[]>> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        if (filter && filter.length > 0) {
            q = q + (q != '' && '&' || '?') + 'filter=' + encodeURI(JSON.stringify(filter))
        }
        http.get('/tests' + q)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
export const PostTests = (_ip: Tests, query?: string): Promise<SimpleAPIResponseWithData<Tests>> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        http.post('/tests' + q, _ip)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
export const GetTestsById = (id: string, query?: string, filter?: TestsSearchCritCondition[]): Promise<SimpleAPIResponseWithData<Tests>> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        if (filter && filter.length > 0) {
            q = q + (q != '' && '&' || '?') + 'filter=' + encodeURI(JSON.stringify(filter))
        }
        http.get(`/tests/${id}` + q)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
export const PutTestsById = (id: string, _ip: Tests, query?: string): Promise<SimpleAPIResponseWithData<Tests>> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        http.put(`/tests/${id}` + q, _ip)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
export const DeleteTestsById = (id: string, query?: string): Promise<SimpleAPIResponse> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        http.delete(`/tests/${id}` + q)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
export const GetTestsById_special = (id: string, query?: string): Promise<SimpleAPIResponse> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        http.get(`/tests/${id}/_special` + q)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
export const GetTests1 = (query?: string, filter?: TestsSearchCritCondition[]): Promise<SimpleAPIResponseWithData<Tests>> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        if (filter && filter.length > 0) {
            q = q + (q != '' && '&' || '?') + 'filter=' + encodeURI(JSON.stringify(filter))
        }
        http.get('/tests1' + q)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
export const GetTests3 = (query?: string): Promise<Tests[]> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        http.get('/tests3' + q)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
