<template>
  <scroller ref="scrollArea" :expand="true" :bottom-margin="20" style="width: 99%">
    <div v-for="(item, index) in items" :key="index">{{ item }}</div>
  </scroller>
</template>

<script setup>
import { useStore } from "vuex";
import { computed, ref, watch } from "vue";
import scroller from "./Scroller.vue";

const store = useStore();
const scrollArea = ref(null);

const levelColor = (level) => {
  if (level == "error") {
    return "text-negative";
  } else {
    return "text-positive";
  }
};

const onLogToggle = function (value) {
  console.log("Hi");
  var cmd;
  if (value) {
    cmd = "start";
  } else {
    cmd = "stop";
  }
  store.dispatch(
    "sendWSCommand",
    JSON.stringify({
      class: "log",
      command: cmd,
    })
  );
  if (!value) {
    store.commit("clearLog");
  }
};

const items = computed({
  get() {
    return store.state.websocket.wslog;
  },
});

watch(
  items,
  (oldValue, newValue) => {
    scrollArea.value.scrollToBottom();
  },
  { deep: true }
);
</script>
