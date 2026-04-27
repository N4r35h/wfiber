// this file was generated DO NOT EDIT
export interface APIResponse<T> {
	success: boolean
	data: T
}
export interface Todo {
	id: number
	title: string
	done: boolean
}
export interface TodoSearchCritCondition {
	field: 'id' | 'title' | 'done'
	condition: 'equals' | 'not equals' | 'contains' | 'not contains' | 'greater than' | 'lesser than'
	value: any
}
import http from './http'
export const GetTodos = (query?: string, filter?: TodoSearchCritCondition[]): Promise<APIResponse<Todo[]>> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        if (filter && filter.length > 0) {
            q = q + (q != '' && '&' || '?') + 'filter=' + encodeURI(JSON.stringify(filter))
        }
        http.get('/todos' + q)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
export const PostTodos = (_ip: Todo, query?: string): Promise<APIResponse<Todo>> => {
    return new Promise((resolve, reject) => {
        let q: string = query || ''
        http.post('/todos' + q, _ip)
            .then(response => {
                return resolve(response.data)
            })
            .catch(reject)
    })
}
