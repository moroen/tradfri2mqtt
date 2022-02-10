<template>
    <div class="q-pa-md" style="max-width: 400px">
        <q-form @submit="onSubmit" @reset="onReset" class="q-gutter-md">
            <q-input
                filled
                v-model="name"
                label="Your name *"
                hint="Name and surname"
                lazy-rules
                :rules="[val => val && val.length > 0 || 'Please type something']"
            />

            <q-input
                filled
                type="number"
                v-model="age"
                label="Your age *"
                lazy-rules
                :rules="[
                    val => val !== null && val !== '' || 'Please type your age',
                    val => val > 0 && val < 100 || 'Please type a real age'
                ]"
            />

            <q-toggle v-model="accept" label="I accept the license and terms" />

            <div>
                <q-btn label="Submit" type="submit" color="primary" />
                <q-btn label="Reset" type="reset" color="primary" flat class="q-ml-sm" />
            </div>
        </q-form>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import { useQuasar } from 'quasar';

const $q = useQuasar()

const name = ref("")
const age = ref(0)
const accept = ref(false)

function onSubmit() {
    console.log("Submit")
    if (accept.value != true) {
        $q.notify({
            color: 'red-5',
            textColor: 'white',
            icon: 'warning',
            message: 'You need to accept the license and terms first'
        })
    }
}


function onReset() {
    console.log("Reset")
    name.value = null
    age.value = null
    accept.value = null
}

</script>
