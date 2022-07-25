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

        <v-autocomplete
            v-model="room"
            label="Room"
            :items="rooms"
            item-title="name"
            item-value="id"
            :rules="[() => !!room || 'This field is required']"
            @update:modelValue="validate"
            required
        >
        </v-autocomplete>


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
        <div v-if="getSelectedRoom()">
          <v-alert v-if="hasOwnOverlaps"
                   density="comfortable"
                   type="error"
                   variant="contained-text"
                   class="mb-2"
          >
            You already have booked a desk in the selected <br/>room for the selected timeframe.
          </v-alert>
          <overlapping-bookings v-else :room="getSelectedRoom()" :date="date" :start="start"
                                :end="end"></overlapping-bookings>

        </div>


      </v-card-text>

      <v-card-actions>
        <v-btn
            :disabled="!valid || hasOwnOverlaps"
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
import {mapGetters, mapState} from "vuex";
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
    end: '17:00',
    hasOwnOverlaps: false
  }),

  mounted() {
    this.$store.dispatch("rooms/fetchRooms")
    this.$store.dispatch("configuration/fetchConfiguration")
  },
  computed: {
    ...mapState('rooms', {rooms: 'rooms', roomsLoading: 'loading'}),
    ...mapState('configuration', {configuration: 'configuration', configurationLoading: 'loading'}),
    ...mapGetters('bookings', {getOverlaps: 'getOverlappingBookings'})
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
      const start = this.getStartDate()
      const end = this.getEndDate()
      const room = this.getSelectedRoom();
      const overlaps = this.getOverlaps(this.getSelectedRoom(), start, end);
      this.hasOwnOverlaps = room !== null ? overlaps.length > 0 : false;

    },
    getSelectedRoom() {
      return this.rooms.find(r => this.room === r.id)
    },
    getStartDate() {
      return moment(this.date).add(this.start)

    },
    getEndDate() {
      return moment(this.date).add(this.end)

    },
    submit() {
      this.validate()
      if (this.valid) {
        const start = this.getStartDate()
        const end = this.getEndDate()
        const room = this.getSelectedRoom()
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