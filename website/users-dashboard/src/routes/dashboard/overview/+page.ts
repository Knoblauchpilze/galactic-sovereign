
import User, { getUser,  } from '$lib/users';

/** @type {import('./$types').PageLoad} */
export async function load({ params }) {
  console.log("params: " + JSON.stringify(params));

  const DUMMY_API_KEY = 'your-key';
  const DUMMY_USER_ID = 'your-id';
  const userResponse = await getUser(DUMMY_API_KEY, DUMMY_USER_ID);

  if (userResponse.error()) {
    // loginError = String(userResponse.details);
  }

  console.log("response: ", JSON.stringify(userResponse));

	return {
    user: new User(userResponse),
	};
}
