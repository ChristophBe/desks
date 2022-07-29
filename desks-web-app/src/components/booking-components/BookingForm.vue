<template>

  <v-card class="card">

    <v-card-title v-if="this.booking">
      Book a Desk
    </v-card-title>
    <v-card-title v-else>
      Edit Desk Booking
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
            :rules="[
                () => !!start || 'This field is required',
                () => startValidationRule() || 'The start time should be in the future'
            ]"
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
          {{ booking? "Save": "Book"}}
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

<script lang="ts">

import {defineComponent} from "vue";
import {mapActions, mapGetters, mapState} from "vuex";
import moment, {Moment} from "moment";
import OverlappingBookings from "@/components/booking-components/OverlappingBookings.vue";
import Room from "@/models/Room";
import Booking from "@/models/Booking";

export default defineComponent({
  name: "BookingForm",
  components: {OverlappingBookings},
  props: {
    booking: {
      type: Object
    }
  },
  data: () => ({
    valid: false,
    date: moment().add(1, 'days').format("YYYY-MM-DD"),
    room: null,
    start: '09:00',
    end: '17:00',
    hasOwnOverlaps: false
  }),

  mounted() {
    this.fetchRooms()
    this.fetchConfiguration()

    if(this.booking != null){
      this.room = this.booking.room.id;
      this.date = moment(this.booking.start).format("YYYY-MM-DD");
      this.start = moment(this.booking.start).format("HH:mm");
      this.end = moment(this.booking.end).format("HH:mm");
      const form = this.$refs.form as HTMLFormElement
      form.validate()
    }
  },
  computed: {
    ...mapState('rooms', {rooms: 'rooms', roomsLoading: 'loading'}),
    ...mapState('configuration', {configuration: 'configuration', configurationLoading: 'loading'}),
    ...mapGetters('bookings', {getOverlaps: 'getOverlappingBookings'})
  },
  methods: {

    ...mapActions("rooms",["fetchRooms"]),
    ...mapActions("configuration",["fetchConfiguration"]),
    ...mapActions("bookings",["saveBooking"]),
    endValidationRule(v:string) {
      const start = moment().add(this.start);
      const end = moment().add(v);
      return end.isAfter(start)
    },
    startValidationRule() {
      return  this.getStartDate().isAfter(moment.now());
    },
    dateMinValidationRule(v:string) {
      return moment(v).isSameOrAfter(moment.now(), "days")
    },
    dateMaxValidationRule(v:string) {
      return moment(v).isSameOrBefore(moment(this.configuration.maximalBookingDate), "days")
    },
    validate() {

      const form = this.$refs.form as HTMLFormElement
      form.validate()
      const start = this.getStartDate()
      const end = this.getEndDate()
      const room = this.getSelectedRoom();
      const overlaps = this.getOverlaps(this.getSelectedRoom(), start, end, this.booking ? [this.booking.id] :[]);
      this.hasOwnOverlaps = room !== null ? overlaps.length > 0 : false;

    },
    getSelectedRoom() {
      return this.rooms.find((r:Room)=> this.room! === r.id)
    },
    getStartDate()  :Moment{
      return moment(this.date).add(this.start)

    },
    getEndDate() :Moment{
      return moment(this.date).add(this.end)

    },
    submit() {
      this.validate()
      if (this.valid) {
        const start = this.getStartDate()
        const end = this.getEndDate()
        const room = this.getSelectedRoom()
        const booking:Partial<Booking> = {
          room:room,
          start: start.toDate(),
          end: end.toDate()
        }

        if(this.booking){
          booking.id = this.booking.id
        }
        this.saveBooking(booking)
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
