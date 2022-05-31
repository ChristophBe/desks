<template>

  <v-card class="card">

    <v-card-title>
      Book a Desk
    </v-card-title>

    <div v-if="roomsLoading || configurationLoading">
      <v-progress-circular
          indeterminate
          color="primary"
      ></v-progress-circular>
    </div>
    <v-form
        v-else
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
            :items="rooms.map(r => r.name)"
            :rules="[() => !!room || 'This field is required']"
            @update:modelValue="validate"

            required
        >
        </v-select>


        <v-text-field
            v-model="date"
            label="Date"
            type="date"
            :rules="[
                () => !!date || 'This field is required',
                (v) => dateMinValidationRule(v) || 'The date can not be in the past',
                (v) => dateMaxValidationRule(v) || 'The date is to far in the future'
             ]"
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
            :rules="[
              () => !!end || 'This field is required',
              (v) => endValidationRule(v) ||  'End has to be after start'
            ]"
            @change="validate"
            required
        ></v-text-field>

        <overlapping-bookings v-if="room" :room="rooms.find(r => r.name === room)" :date="date" :start="start"
                              :end="end"></overlapping-bookings>
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
import OverlappingBookings from "@/components/OverlappingBookings";


export default defineComponent({
  name: "AddBookingForm",
  components: {OverlappingBookings},
  data: () => ({
    valid: false,
    date: moment().add(1, 'days').format("YYYY-MM-DD"),
    room: null,
    start: '09:00',
    end: '17:00'
  }),

  mounted() {
    this.$store.dispatch("rooms/fetchRooms")
    this.$store.dispatch("configuration/fetchConfiguration")
  },
  computed: {
    ...mapState('rooms', {rooms: 'rooms', roomsLoading: 'loading'}),
    ...mapState('configuration', {configuration: 'configuration', configurationLoading: 'loading'}),
  },
  methods: {
    endValidationRule(v) {
      const start = moment().add(this.start);
      const end = moment().add(v);
      return end.isAfter(start)
    },
    dateMinValidationRule(v) {
      return moment(v).isSameOrAfter(moment.now(), "days")
    },
    dateMaxValidationRule(v) {
      return moment(v).isSameOrBefore(moment(this.configuration.maximalBookingDate), "days")
    },
    validate() {
      this.$refs.form.validate()
    },
    submit() {
      this.validate()
      if (this.valid) {
        const start = moment(this.date).add(this.start)
        const end = moment(this.date).add(this.end)
        const room = this.rooms.find(r => this.room === r.name)
        const booking = {
          room,
          start,
          end
        }

        this.$store.dispatch('bookings/bookDesk', booking)
        this.$emit("close")
      }

    }
  }
})
</script>

<style scoped>
.card {
  min-width: 400px;
}
</style>