
/** @type {import('./$types').PageLoad} */
export function load({ params }) {
  console.log("params: " + JSON.stringify(params));

	return {
		id: "00-00-00-01",
    email: "random@e.mail",
    password: "juice",
    createdAt: new Date()
	};
}
