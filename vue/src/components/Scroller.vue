<template>
  <div ref="scrollContainer">
    <q-scroll-area
      ref="scrollArea"
      class="col"
      :visible="true"
      v-bind:style="{ height: scrollHeight }"
    >
      <slot></slot>
    </q-scroll-area>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { stringifyQuery } from "vue-router";

const scrollContainer = ref(null);
const scrollArea = ref(null);

const scrollHeight = ref(null);

const props = defineProps({
  expand: {
    type: Boolean,
    default: false,
  },
  bottomMargin: {
    type: Number,
    default: 20,
  },
  height: {
    type: String,
    default: "300px",
  },
});

const heightFn = () => {
  let h =
    window.innerHeight -
    scrollContainer.value.getBoundingClientRect().top -
    props.bottomMargin;
  return h + "px";
};

const scrollToBottom = () => {
  const scrollTarget = scrollArea.value.getScrollTarget();
  scrollArea.value.setScrollPosition("vertical", scrollTarget.scrollHeight, 300);
};

const resize = () => {
  scrollHeight.value = heightFn();
};

onMounted(() => {
  if (props.expand) {
    scrollHeight.value = heightFn();
    window.addEventListener("resize", resize);
  } else {
    scrollHeight.value = props.height;
  }
});

onUnmounted(() => {
  if (props.expand) {
    window.removeEventListener("resize", resize);
  }
});

defineExpose({ scrollToBottom });
</script>
