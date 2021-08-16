/*
export function someMutation (state) {
}
*/

export const updateTradfri = (state, payload) => {
    state[payload.section][payload.key] = payload.value
  }