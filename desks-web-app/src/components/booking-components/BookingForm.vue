<template>

  <v-card class="card">

    <v-card-title v-if="booking && booking.id">
      Edit Booking
    </v-card-title>
    <v-card-title v-else>
      Book a Desk
    </v-card-title>

    <v-form
        ref="form"
        v-model="valid"
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
                   variant="text"
                   class="mb-2"
          >
            You already have booked a desk in the selected <br/>room for the selected timeframe.
          </v-alert>
          <overlapping-bookings v-else
                                :room="getSelectedRoom()"
                                :date="date"
                                :start="start"
                                :end="end"></overlapping-bookings>
        </div>
      </v-card-text>

      <v-card-actions>
        <v-btn
            :disabled="!valid || hasOwnOverlaps"
            @click="submit"
        >
          {{ booking && booking.id ? "Save" : "Book" }}
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
    },
    defaults: {
      type: Object
    }
  },
  data: () => ({
    valid: true,
    date: moment().add(1, 'days').format("YYYY-MM-DD"),
    room: null,
    start: '09:00',
    end: '17:00',
    hasOwnOverlaps: false
  }),

  mounted() {

    this.fetchRooms();
    this.fetchBookingDefaults();
    this.fetchConfiguration();
    if (this.booking != null) {
      this.room = this.booking.room.id;
      this.date = moment(this.booking.start).format("YYYY-MM-DD");
      this.start = moment(this.booking.start).format("HH:mm");
      this.end = moment(this.booking.end).format("HH:mm");

    } else {
      this.setDefaults()

    }
    this.validate()
  },
  computed: {
    ...mapState('rooms', {rooms: 'rooms'}),
    ...mapState('configuration', {configuration: 'configuration'}),
    ...mapGetters('bookings', {getOverlaps: 'getOverlappingBookings'})
  },

  methods: {
    ...mapActions("rooms", ["fetchRooms"]),
    ...mapActions("configuration", ["fetchConfiguration"]),
    ...mapActions("defaults", ["fetchBookingDefaults"]),
    ...mapActions("bookings", ["saveBooking"]),

    setDefaults() {
      if (this.booking == null && this.defaults) {
        this.date = moment(this.defaults.start).format("YYYY-MM-DD");
        this.start = moment(this.defaults.start).format("HH:mm");
        this.end = moment(this.defaults.end).format("HH:mm");

        if (this.defaults.room) {
          this.room = this.defaults.room.id
        }
      }
    },
    endValidationRule(v: string) {
      const start = moment().add(this.start);
      const end = moment().add(v);
      return end.isAfter(start)
    },
    startValidationRule() {
      return this.getStartDate().isAfter(moment.now());
    },
    dateMinValidationRule(v: string) {
      return moment(v).isSameOrAfter(moment.now(), "days")
    },
    dateMaxValidationRule(v: string) {
      return moment(v).isSameOrBefore(moment(this.configuration.maximalBookingDate), "days")
    },
    validate() {
      const form = this.$refs.form as HTMLFormElement
      if (form) {
        form.validate()
      }

      const start = this.getStartDate()
      const end = this.getEndDate()
      const room = this.getSelectedRoom();
      const overlaps = this.getOverlaps(this.getSelectedRoom(), start, end, this.booking && this.booking.id ? [this.booking.id] : []);
      this.hasOwnOverlaps = room !== null ? overlaps.length > 0 : false;
    },
    getSelectedRoom() {
      return this.rooms.find((r: Room) => this.room! === r.id)
    },
    getStartDate(): Moment {
      return moment(this.date).add(this.start)

    },
    getEndDate(): Moment {
      return moment(this.date).add(this.end)

    },
    submit() {
      this.validate()
      if (this.valid) {
        const start = this.getStartDate()
        const end = this.getEndDate()
        const room = this.getSelectedRoom()
        const booking: Partial<Booking> = {
          room: room,
          start: start.toDate(),
          end: end.toDate()
        }

        if (this.booking) {
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
