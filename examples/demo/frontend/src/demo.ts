import { PostTodos, type Todo } from "./api/api.gen"

const payload: Todo = {
	id: 42,
	title: "Record short wfiber demo",
	done: false,
}

const doubledID = payload.id * 2
console.log("Doubled ID:", doubledID)

async function run() {
	const created = await PostTodos(payload)
	console.log("Created todo:", created)
}

void run()
