<script setup>
import { ref, computed } from "vue";
import { useMapStore } from "../../store/mapStore";

const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);
// const mapStore = useMapStore();

const targetArea = ref(null);
// const AreaColor = ref(props.chart_config.color[0]);
const mousePosition = ref({ x: null, y: null });
const refMousePosition = ref({ x: null, y: null });
const selectedIndex = ref(null);
const mapStore = useMapStore();

// Parse Venn Diagram Data
const VennData = computed(() => {
	let readData = {};
	let count = {};
	let category = [];
	let sum = 0;

	let i = 0;
	while (props.series[i].name !== null) {
		i++;
	}
	i++;

	let len = props.series[i].data.length;
	if (len === 2) {
		for (let i = 0; i < 3; i++) {
			count[i.toString()] = 0;
		}
	} else if (len === 3) {
		for (let i = 0; i < 7; i++) {
			count[i.toString()] = 0;
		}
	}

	for (let k = 0; k < len; k++) {
		category.push(props.series[i].data[k]);
	}
	i++;

	while (i < props.series.length) {
		sum++;
		const item = props.series[i];
		const dat = item.data;
		readData[item.name] = dat;
		if (len === 2) {
			if ((dat[0] === "1") & (dat[1] === "1")) {
				count["2"]++;
			} else if ((dat[0] === "1") & (dat[1] === "0")) {
				count["0"]++;
			} else if ((dat[0] === "0") & (dat[1] === "1")) {
				count["1"]++;
			}
		} else if (len === 3) {
			if ((dat[0] === "1") & (dat[1] === "1") & (dat[2] === "1")) {
				count["6"]++;
			} else if ((dat[0] === "1") & (dat[1] === "0") & (dat[2] === "0")) {
				count["0"]++;
			} else if ((dat[0] === "0") & (dat[1] === "1") & (dat[2] === "0")) {
				count["1"]++;
			} else if ((dat[0] === "0") & (dat[1] === "0") & (dat[2] === "1")) {
				count["2"]++;
			} else if ((dat[0] === "1") & (dat[1] === "1") & (dat[2] === "0")) {
				count["3"]++;
			} else if ((dat[0] === "1") & (dat[1] === "0") & (dat[2] === "1")) {
				count["4"]++;
			} else if ((dat[0] === "0") & (dat[1] === "1") & (dat[2] === "1")) {
				count["5"]++;
			}
		}
		i++;
	}

	const output = {
		sum: sum,
		len: len,
		readData: readData,
		count: count,
		category: category,
	};
	return output;
});

const tooltipPosition = computed(() => {
	return {
		left: `${mousePosition.value.x - 10}px`,
		top: `${mousePosition.value.y - 54}px`,
	};
});
function toggleActive(e) {
	targetArea.value = e.target.dataset.name;
}
function toggleActiveToNull() {
	targetArea.value = null;
}
function updateMouseLocation(e) {
	const container = document.getElementsByClassName("svg");
	const rect = container[0].getBoundingClientRect();
	refMousePosition.value.x =
		((e.pageX - rect.left) / (rect.right - rect.left)) * 550;
	refMousePosition.value.y =
		((e.pageY - rect.top) / (rect.bottom - rect.top)) * 550;
	mousePosition.value.x = e.pageX;
	mousePosition.value.y = e.pageY;
}

function handleDataSelection(index) {
	if (!props.chart_config.map_filter) {
		return;
	}
	if (index !== selectedIndex.value) {
		mapStore.addLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`,
			props.chart_config.map_filter[2],
			props.chart_config.map_filter[3][index]
		);
		selectedIndex.value = index;
	} else {
		mapStore.clearLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`
		);
		selectedIndex.value = null;
	}
}
</script>

<template>
	<div v-if="activeChart === 'VennDiagram'" class="venndiagram">
		<div class="venndiagram-title">
			<h5>總機構數量</h5>
			<h6>{{ VennData.sum }} {{ chart_config.unit }}</h6>
		</div>
		<div class="venndiagram-chart">
			<div v-if="VennData.len === 2" class="circle">
				<svg
					viewBox="0 0 650 650"
					class="svg"
					xmlns="http://www.w3.org/2000/svg"
				>
					<g transform="translate(50 50)">
						<path
							data-name="1"
							fill="#4574a1"
							stroke="#4f4f4f"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '1' || selectedIndex === 1,
								'initial-animation-3': true,
							}"
							d="M 275 145.09618943233
							A 150 150 0 1 1 275 404.90381056767
							A 150 150 0 0 0 275 145.09618943233
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(1)"
						/>
						<path
							data-name="0"
							fill="#9c54ab"
							stroke="#4f4f4f"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '0' || selectedIndex === 0,
								'initial-animation-2': true,
							}"
							d="M 275 145.09618943233
							A 150 150 0 1 0 275 404.90381056767
							A 150 150 0 0 1 275 145.09618943233
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(0)"
						/>
						<path
							data-name="2"
							fill="#7164a6"
							stroke="#4f4f4f"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '2' || selectedIndex === 2,
								'initial-animation-1': true,
							}"
							d="M 275 145.09618943233
							A 150 150 0 0 0 275 404.90381056767
							A 150 150 0 0 0 275 145.09618943233
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(2)"
						/>
					</g>
				</svg>
				<div class="text">
					<div class="text-2 text-3">
						<div class="square-1"></div>
						<p>&nbsp;&nbsp;</p>
						<h6 class="text-color">{{ VennData.category[0] }}</h6>
					</div>
					<div class="text-2 text-3">
						<div class="square-2"></div>
						<p>&nbsp;&nbsp;</p>
						<h6 class="text-color">{{ VennData.category[1] }}</h6>
					</div>
				</div>
			</div>
			<div v-if="VennData.len === 3" class="circle">
				<svg
					viewBox="0 0 650 650"
					class="svg"
					xmlns="http://www.w3.org/2000/svg"
					@mousemove="updateMouseLocation"
					@mouseleave="toggleActiveToNull"
				>
					<g transform="translate(50 50)">
						<path
							data-name="0"
							fill="#519a98"
							stroke="#828282"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '0' || selectedIndex === 0,
								'initial-animation-5': true,
							}"
							d="M 125.6327126328 213.76275643042
							A 150 150 0 0 1 275 227.52551286084
							A 150 150 0 0 1 424.3672873672 213.76275643042
							A 150 150 0 1 0 125.6327126328 213.76275643042
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(0)"
						/>
						<path
							data-name="1"
							fill="#93bb77"
							stroke="#828282"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '1' || selectedIndex === 1,
								'initial-animation-6': true,
							}"
							d="M 125.6327126328 213.76275643042
							A 150 150 0 0 0 212.23525301124 336.23724356958
							A 150 150 0 0 0 275 472.47448713916
							A 150 150 0 1 1 125.6327126328 213.76275643042
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(1)"
						/>
						<path
							data-name="2"
							fill="#D6A059"
							stroke="#828282"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '2' || selectedIndex === 2,
								'initial-animation-7': true,
							}"
							d="M 424.3672873672 213.76275643042
							A 150 150 0 0 1 337.76474698876 336.23724356958
							A 150 150 0 0 1 275 472.47448713916
							A 150 150 0 1 0 424.3672873672 213.76275643042
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(2)"
						/>

						<path
							data-name="3"
							fill="#72ab88"
							stroke="#828282"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '3' || selectedIndex === 3,
								'initial-animation-2': true,
							}"
							d="M 125.6327126328 213.76275643042
							A 150 150 0 0 1 275 227.52551286084
							A 150 150 0 0 0 212.23525301124 336.23724356958
							A 150 150 0 0 1 125.6327126328 213.76275643042
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(3)"
						/>
						<path
							data-name="4"
							fill="#949d79"
							stroke="#828282"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '4' || selectedIndex === 4,
								'initial-animation-4': true,
							}"
							d="M 275 227.52551286084
							A 150 150 0 0 1 424.3672873672 213.76275643042
							A 150 150 0 0 1 337.76474698876 336.23724356958
							A 150 150 0 0 0 275 227.52551286084
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(4)"
						/>
						<path
							data-name="5"
							fill="#b5b068"
							stroke="#828282"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '5' || selectedIndex === 5,
								'initial-animation-3': true,
							}"
							d="M 212.23525301124 336.23724356958
							A 150 150 0 0 0 337.76474698876 336.23724356958
							A 150 150 0 0 1 275 472.47448713916
							A 150 150 0 0 1 212.23525301124 336.23724356958
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(5)"
						/>
						<path
							data-name="6"
							fill="#93a678"
							stroke="#828282"
							stroke-width="2"
							:class="{
								'active-block':
									targetArea === '6' || selectedIndex === 6,
								'initial-animation-1': true,
							}"
							d="M 275 227.52551286084
							A 150 150 0 0 0 212.23525301124 336.23724356958
							A 150 150 0 0 0 337.76474698876 336.23724356958
							A 150 150 0 0 0 275 227.52551286084
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(6)"
						/>
					</g>
				</svg>
				<div class="text">
					<div class="text-1">
						<div class="text-2">
							<div class="square-3"></div>
							<p>&nbsp;&nbsp;</p>
							<h6 class="text-color">
								{{ VennData.category[0] }}
							</h6>
						</div>
						<div class="text-2">
							<div class="square-4"></div>
							<p>&nbsp;&nbsp;</p>
							<h6 class="text-color">
								{{ VennData.category[1] }}
							</h6>
						</div>
					</div>
					<div class="text-2 text-3">
						<div class="square-5"></div>
						<p>&nbsp;&nbsp;</p>
						<h6 class="text-color">{{ VennData.category[2] }}</h6>
					</div>
				</div>
			</div>
			<Teleport to="body">
				<!-- The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css -->
				<div
					v-if="targetArea"
					class="venndiagram-chart-info chart-tooltip"
					:style="tooltipPosition"
				>
					<h6>數量</h6>
					<span
						>{{ VennData.count[targetArea] }}
						{{ chart_config.unit }}</span
					>
				</div>
			</Teleport>
		</div>
	</div>
</template>

<style scoped lang="scss">
.venndiagram {
	max-height: 100%;
	position: relative;
	overflow-y: scroll;

	&-title {
		display: flex;
		flex-direction: column;
		justify-content: center;
		position: absolute;
		left: 0;
		top: 0;
		margin: 0.5rem 0 -0.5rem;

		h5 {
			color: var(--color-complement-text);
		}

		h6 {
			color: var(--color-complement-text);
			font-size: var(--font-m);
			font-weight: 400;
		}

		&-legend {
			display: flex;
			justify-content: space-between;

			div {
				position: relative;
				width: 3rem;
				margin: 0 4px;
				border-radius: 5px;
			}

			div:before {
				content: "";
				width: 3rem;
				height: var(--font-l);
				position: absolute;
				top: 0;
				left: 0;
				background: linear-gradient(
					270deg,
					rgba(40, 42, 44, 1),
					rgba(40, 42, 44, 0.2)
				);
			}

			p {
				color: var(--color-complement-text);
			}
		}
	}

	&-chart {
		display: flex;
		justify-content: center;

		svg {
			width: 80%;
			-webkit-transform: translateX(10%);
			-ms-transform: translateX(10%);
			transform: translateX(10%);
		}

		&-info {
			position: fixed;
			z-index: 20;
		}
	}
}

.text {
	display: flex;
	flex-direction: column;
	transform: translateY(-20px);
}

.text-1 {
	display: flex;
	flex-direction: row;
	justify-content: space-evenly;
}
.text-2 {
	display: flex;
	flex-direction: row;
}

.text-3 {
	transform: translateX(25%);
}

.text-color {
	color: #828282;
}

.square-1 {
	width: 10px;
	height: 10px;
	border-radius: 3px;
	background-color: #9450a2;
	transform: translateY(3px);
}

.square-2 {
	width: 10px;
	height: 10px;
	border-radius: 3px;
	background-color: #44729e;
	transform: translateY(3px);
}

.square-3 {
	width: 10px;
	height: 10px;
	border-radius: 3px;
	background-color: #519a98;
	transform: translateY(3px);
}

.square-4 {
	width: 10px;
	height: 10px;
	border-radius: 3px;
	background-color: #93bb77;
	transform: translateY(3px);
}

.square-5 {
	width: 10px;
	height: 10px;
	border-radius: 3px;
	background-color: #d6a059;
	transform: translateY(3px);
}

@keyframes ease-in {
	0% {
		opacity: 0;
	}

	100% {
		opacity: 1;
	}
}

@for $i from 1 to 10 {
	.initial-animation-#{$i} {
		animation-name: ease-in;
		animation-duration: 0.3s;
		animation-delay: 0.2s * ($i - 1);
		animation-timing-function: linear;
		animation-fill-mode: forwards;
		opacity: 0;
	}
}
.active-block {
	transform: translateY(-7px);
}
</style>
