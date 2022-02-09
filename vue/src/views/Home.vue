<template>
    <div class="q-pa-md" style="max-width: 300px">
        <form @submit.prevent.stop="onSubmit" @reset.prevent.stop="onReset" class="q-gutter-md">
            <q-input
                
                filled
                v-model="name"
                label="Your name *"
                hint="Name and surname"
                :rules="nameRules"
            />
        </form>
    </div>
</template>


<script>
import { useQuasar } from 'quasar'
import { computed, ref } from "vue";
import { useStore } from 'vuex';

export default {
    
    setup() {
        const store = useStore()

        const name = computed({
            get: () => store.state.text,
            set: (val) => {
                store.commit("setText", val);
            }
        })

        return {
            name,
            nameRules: [
                val => (val && val.length > 0) || 'Please type something'
            ],
        }
    }
}
</script>
