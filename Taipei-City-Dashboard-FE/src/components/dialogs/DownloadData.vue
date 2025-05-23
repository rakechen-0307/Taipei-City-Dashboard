<!-- eslint-disable indent -->
<!-- eslint-disable no-mixed-spaces-and-tabs -->
<!-- Developed by Taipei Urban Intelligence Center 2023-2024-->

<script setup>
import { ref, computed } from "vue";
import { useDialogStore } from "../../store/dialogStore";

import { jsonToCsv } from "../../assets/utilityFunctions/jsonToCsv";
import DialogContainer from "./DialogContainer.vue";

const props = defineProps(["content"]);

const dialogStore = useDialogStore();
const content = computed(() => props.content || dialogStore.moreInfoContent);

// Stores the inputted dashboard name
const name = ref(content.value.name);
// Stores the file type
const fileType = ref("JSON");

const parsedJson = computed(() => {
	let json = {};
	json.data = content.value.chart_data;
	if (content.value.chart_config.categories) {
		json.categories = content.value.chart_config.categories;
	}

	const jsonString = encodeURIComponent(JSON.stringify(json));
	// const base64Json = btoa(jsonString)
	return jsonString;
});

const parsedCsv = computed(() => {
	const csvString = content.value.chart_data
		? jsonToCsv(
				content.value.chart_data,
				content.value.chart_config
		  )
		: "";
	return encodeURI(csvString);
});

function handleSubmit() {
	handleClose();
}
function handleClose() {
	name.value = content.value.name;
	dialogStore.dialogs.downloadData = false;
}
</script>

<template>
  <DialogContainer
    :dialog="`downloadData`"
    @on-close="handleClose"
  >
    <div class="downloaddata">
      <h2>下載資料</h2>
      <div class="downloaddata-input">
        <h3>請輸入檔名</h3>
        <input
          v-model="name"
          type="text"
          :minlength="1"
          required
        >
      </div>
      <h3>請選擇檔案格式</h3>
      <div>
        <input
          id="JSON"
          v-model="fileType"
          class="downloaddata-radio"
          type="radio"
          value="JSON"
        >
        <label for="JSON">
          <div />
          JSON
        </label>
        <input
          id="CSV"
          v-model="fileType"
          class="downloaddata-radio"
          type="radio"
          value="CSV"
        >
        <label for="CSV">
          <div />
          CSV (UTF-8)
        </label>
      </div>
      <div class="downloaddata-control">
        <button
          class="downloaddata-control-cancel"
          @click="handleClose"
        >
          取消
        </button>
        <button
          v-if="name && fileType === 'JSON'"
          class="downloaddata-control-confirm"
          @click="handleSubmit"
        >
          <a
            :href="`data:application/json;charset=utf-8,${parsedJson}`"
            :download="`${name}.json`"
          >下載JSON</a>
        </button>
        <button
          v-if="name && fileType === 'CSV'"
          class="downloaddata-control-confirm"
          @click="handleSubmit"
        >
          <a
            :href="`data:text/csv;charset=utf-8,${parsedCsv}`"
            :download="`${name}.csv`"
          >下載CSV</a>
        </button>
      </div>
    </div>
  </DialogContainer>
</template>

<style scoped lang="scss">
.downloaddata {
	width: 300px;

	h3 {
		margin-bottom: 0.5rem;
		font-size: var(--font-s);
		font-weight: 400;
		color: var(--color-complement-text);
	}

	&-input {
		display: flex;
		flex-direction: column;
		margin: 0.5rem 0;
	}

	&-radio {
		display: none;

		&:checked + label {
			color: white;

			div {
				background-color: var(--color-highlight);
			}
		}

		&:hover + label {
			color: var(--color-highlight);

			div {
				border-color: var(--color-highlight);
			}
		}
	}

	label {
		position: relative;
		display: flex;
		align-items: center;
		margin-bottom: 2px;
		font-size: var(--font-s);
		color: var(--color-complement-text);
		transition: color 0.2s;
		cursor: pointer;

		div {
			width: calc(var(--font-s) / 2);
			height: calc(var(--font-s) / 2);
			margin-right: 4px;
			padding: calc(var(--font-s) / 4);
			border-radius: 50%;
			border: 1px solid var(--color-border);
			transition: background-color 0.2s;
		}
	}

	&-control {
		display: flex;
		justify-content: flex-end;

		&-cancel {
			margin: 0 2px;
			padding: 4px 6px;
			border-radius: 5px;
			transition: color 0.2s;

			&:hover {
				color: var(--color-highlight);
			}
		}

		&-confirm {
			margin: 0 2px;
			padding: 4px 10px;
			border-radius: 5px;
			background-color: var(--color-highlight);
			transition: opacity 0.2s;

			&:hover {
				opacity: 0.8;
			}
		}
	}
}
</style>
