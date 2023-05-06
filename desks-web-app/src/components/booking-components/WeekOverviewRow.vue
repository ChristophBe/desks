<template>


  <v-row class="mt-4">
    <v-col>
      <h1 class="text-h4">Desk Availability</h1>
      <h3 class="text-subtitle-1">{{ $format.date(dateRanges[window][0]) }} - {{ $format.date(dateRanges[window][4]) }} | {{ bookingDefaults?.room?.name }}</h3>
    </v-col>

    <v-col align-self="center" class="d-flex justify-end flex-grow-0" :style="{width: 'fit-content'}">
      <v-btn icon variant="text"
             @click="()=>previousDateRange()"
             :disabled="isPreviousWeekDisabled()">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-btn icon variant="text" @click="()=>nextDateRange()">
        <v-icon>mdi-arrow-right</v-icon>
      </v-btn>
    </v-col>
  </v-row>
  <v-row>
    <v-col>
      <v-window v-model="window">
        <v-window-item v-for="(scope,n) in dateRanges" :value="n" :key="n">
          <week-overview-window :dateScope="scope" @add-booking="(booking) => $emit('book', booking)"></week-overview-window>
        </v-window-item>
      </v-window>
    </v-col>

  </v-row>

</template>
<script lang="ts">
import {defineComponent} from "vue";
import moment from "moment/moment";
import {Moment} from "moment";
import WeekOverviewWindow from "@/components/booking-components/WeekOverviewWindow.vue";
import {mapState} from "vuex";

interface DeskAvailabilityState {
  dateRanges: Moment[][]
  window: number
}

export default defineComponent({
  name: 'desk-availability',
  components: {WeekOverviewWindow},
  data: (): DeskAvailabilityState => ({
    window: 0,
    dateRanges: []
  }),
  computed:{
    ...mapState("defaults",["bookingDefaults"])
  },
  created() {
    this.addDateRange(moment().startOf("day")) // init first range
    this.addNextDateRange() // add next range
  },
  methods: {
    addNextDateRange() {
      const additionalWindow = (this.dateRanges[this.window][4]).add(1, "day").startOf("day")
      this.addDateRange(additionalWindow)
    },
    nextDateRange() {
        this.addNextDateRange();
      this.window++
    },
    addDateRange(firstDate: Moment) {
      this.dateRanges.push(this.getNext5BusinessDays(firstDate));
    },
    previousDateRange() {
      this.window--
    },
    isPreviousWeekDisabled(): boolean {
      return this.window <= 0
    },

    isWeekend(day: Moment) {
      return day.isoWeekday() > 5;
    },

    getNext5BusinessDays(start: moment.Moment): Moment[] {
      const days = new Array<Moment>(7);
      for(var i = 0; i < 7; i++) days[i] = moment(start).add(i, "day")
      return days.filter(d => !this.isWeekend(d));
    }
  }
});
</script>
