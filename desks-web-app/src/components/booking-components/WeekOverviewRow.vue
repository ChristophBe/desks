<template>


  <v-row class="mt-4">
    <v-col>
      <h1 class="text-h4">{{ isPreviousWeekDisabled() ? "This Week" : "Week Overview" }}</h1>
      <h3 class="text-subtitle-1">{{$format.date(calculateNthDayOfWeek(1))}} - {{$format.date(calculateNthDayOfWeek(5))}} | {{bookingDefaults?.room?.name}}</h3>
    </v-col>

    <v-col align-self="center" class="d-flex justify-end">
      <v-btn icon variant="text"
             @click="()=>previousWeek()"
             :disabled="isPreviousWeekDisabled()">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>
      <v-btn icon variant="text" @click="()=>nextWeek()">
        <v-icon>mdi-arrow-right</v-icon>
      </v-btn>
    </v-col>
  </v-row>
  <v-row>
    <v-col>
      <v-window v-model="window">
        <v-window-item v-for="(week,n) in weeks" :value="n" :key="n">
          <week-overview-window :start-of-week="week"  @add-booking="(booking) => $emit('book', booking)"></week-overview-window>
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
  startOfWeek: Moment
  weeks: Moment[]
  window: number
}

export default defineComponent({
  name: 'desk-availability',
  components: {WeekOverviewWindow},
  data: (): DeskAvailabilityState => ({
    window: 0,
    weeks: [moment().startOf("week"), moment().startOf("week").add(1, "week")],
    startOfWeek: moment().startOf("week")
  }),

  computed:{
    ...mapState("defaults",["bookingDefaults"])
  },
  methods: {
    nextWeek() {
      if(this.window == this.weeks.length - 2){
        this.weeks = [...this.weeks,moment(this.weeks[this.weeks.length-1]).add(1, "week")]
      }
      this.window++

      this.startOfWeek = this.weeks[this.window]
    },
    previousWeek() {
      this.window--
      this.startOfWeek = this.weeks[this.window]

    },
    isPreviousWeekDisabled(): boolean {
      return this.window <= 0
    },
    isDisabled: (startOfDay: Moment): boolean => startOfDay.isBefore(moment().startOf("day")),

    calculateNthDayOfWeek(n: number): Moment {
      return moment(this.startOfWeek).add(n, "day")
    },

  }
});
</script>