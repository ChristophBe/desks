<template>

  <v-row>
    <template v-for="businessDay in dateScope" :key="businessDay">
      <v-col class="d-flex flex-column justify-end" :class="{days: true, monday: isMonday(businessDay)}" v-if="this.bookingDefaults">
        <span id="dayNote" v-if="businessDay.startOf('day').isSame(today.startOf('day'), 'day')">Today</span>
        <span id="dayNote" v-else>Week {{businessDay.isoWeek()}}</span>
        <availability-card
            v-if="this.bookingDefaults"
            :startOfDay="businessDay"
            :room="this.bookingDefaults.room"
            @add-booking="bookForDay(businessDay)"
        ></availability-card>
      </v-col>
    </template>

  </v-row>

</template>
<style>

.days {
  padding-block: 1em;
}

.monday {
  border-left: 1px solid rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity));
}

.days:first-child {
  border-left: 0 solid rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity));
}


#dayNote {
  display: none;
  margin-bottom: 1em;
  color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity));
  line-height: 1;
  font-size: 1em;
}

.monday #dayNote, .days:first-child #dayNote {
  display: block;
}

</style>
<script lang="ts">
import {defineComponent, PropType} from "vue";
import {mapActions, mapGetters, mapState} from "vuex";
import moment from "moment/moment";
import BookingUtils from "@/utils/booking-utils";
import Booking from "@/models/Booking";
import {Moment, MomentInput} from "moment";
import AvailabilityCard from "@/components/booking-components/AvailabilityCard.vue"

export default defineComponent({
  name: 'weekOverviewWindow',
  props:{
    dateScope: {
      type: Object as PropType<Moment[]>,
      required: true
    }
  },
  components: {AvailabilityCard},

  computed: {
    ...mapState("defaults", ["bookingDefaults", "bookingDefaultsLoading"]),
    ...mapGetters("bookings", ["getBookingsByRoomAndDay"]),
    ...mapState("bookings", ["bookings"]),
    ...mapState('user', ['user']),
    today: () => moment(),
  },

  beforeMount() {
    this.fetchBookingDefaults()
    this.fetchBookings()
  },
  watch: {
    async bookingDefaults() {
      await this.fetchBookings()
    },
  },
  methods: {

    ...mapActions("defaults", ["fetchBookingDefaults"]),
    ...mapActions("bookings", ["fetchBookingsForRange"]),

    isDisabled: (startOfDay: Moment): boolean => startOfDay.isBefore(moment().startOf("day")),

    isMonday(day: Moment) {
      return day.isoWeekday() === 1;
    },

    setDay(base:MomentInput, day: Moment): Moment{
      const baseStartOfDay = moment(base).startOf("day")
      const secondsFromStartOfDay = moment(base).diff(baseStartOfDay, "seconds")
      return moment(day).startOf("day").add(secondsFromStartOfDay, "seconds");
    },

    bookForDay(startOfDay: Moment) {
        console.log("startOfDay ", startOfDay.toISOString())
      let start = this.setDay(this.bookingDefaults.start, startOfDay);
      start = start.isBefore(moment.now()) ? BookingUtils.roundToNextFiveMinutes() : start
        console.log("start ", start.toISOString())
      let end = this.setDay(this.bookingDefaults.end, startOfDay)
      end = end.isSameOrAfter(start) ? end : moment(start).add(1, "hour")
        console.log("end ", end.toISOString())
      const booking: Partial<Booking> = {
        room: this.bookingDefaults.room,
        start: start.toDate(),
        end: end.toDate()
      }
      this.$emit("add-booking", booking)
    },

    async fetchBookings() {
      if (this.bookingDefaults && this.bookingDefaults.room) {
        await this.fetchBookingsForRange({
            roomId: this.bookingDefaults.room.id,
            from: this.dateScope[0],
            to: moment(this.dateScope[this.dateScope.length - 1]).endOf("day")
        })
      }
    },
  }
});
</script>
