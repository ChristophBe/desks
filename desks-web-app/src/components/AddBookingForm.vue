<template>

  <v-card class="card">

      <v-card-title>
        Book a Desk
      </v-card-title>
      <v-form
          ref="form"
          v-model="valid"
          @submit.prevent="submit"
          @change="validate"
          lazy-validation
      >
        <v-card-text>
          <v-select
              v-model="room"
              label="Room"
              :items="rooms.map(room => room.name)"
              :rules="[() => !!room || 'This field is required']"
              required
          >
          </v-select>
          <v-text-field
              v-model="date"
              label="Date"
              type="date"
              :rules="[() => !!date || 'This field is required']"
              @change="validate"
              required
          ></v-text-field>
          <v-text-field
              v-model="start"
              label="Start"
              type="time"
              :rules="[() => !!start || 'This field is required']"
              @change="validate"
              required
          ></v-text-field>
          <v-text-field
              v-model="end"
              label="End"
              type="time"
              :rules="[() => !!end || 'This field is required']"
              @change="validate"
              required
          ></v-text-field>


        </v-card-text>

        <v-card-actions>
          <v-btn
              :disabled="!valid"
              @click="submit"
          >
            Book
          </v-btn>
          <v-btn
              @click="$emit('close')"
          >
            Cancel
          </v-btn>

        </v-card-actions>

      </v-form>
    </v-card>
</template>

<script>

import {defineComponent} from "vue";
import {mapState} from "vuex";
import moment from "moment";


export default defineComponent({
  name:"AddBookingForm",
  data: ()=>({
    valid:false,
    date: null,
    room: null,
    start: '09:00',
    end:'17:00',
  }),

  mounted() {
    this.$store.dispatch("fetchRooms")
  },
  computed: mapState(['rooms', 'user']),
  methods: {
    validate() {
      this.$refs.form.validate()
    },
    submit () {
      this.validate()
      if(this.valid){
        const start = moment(this.date).add(this.start)
        const end = moment(this.date).add(this.end)
        const room = this.rooms.find(r => this.room === r.name)
        const booking = {
          user: this.user,
          room,
          start,
          end
        }

        this.$store.dispatch('bookDesk', booking)
        this.$emit("close")
      }

    }
  }
})
</script>

<style scoped>
.card{
  min-width: 400px;
}
</style>