import { ref } from "vue";

export default function () {
  return {
    tradfri: {
      gateway: "127.0.1.1",
      ident: ref(null),
      key: ref(null),
      keepAlive: 0
    }
    //
  }
}
