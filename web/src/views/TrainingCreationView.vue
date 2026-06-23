<script setup lang="ts">
import type { Question, RadioQuestion, Training } from '@/models/trainings'
import { ref } from 'vue'

const training = ref<Training>({
  id: 0,
  title: 'New Training',
  description: 'Gets to creating!',
  questions: [],
})
const answers = ref<Record<string, string | number | boolean>>({})

function addQuestion() {
  training.value.questions.push({
    id: 'q1',
    label: 'New Question',
    required: true,
    type: 'radio',
    options: ['Answer'],
  })
}

function changeQuestionToDifferentType(question: Question, newType: Question['type']): Question {
  const base = {
    id: question.id,
    label: question.label,
    required: question.required,
  }

  switch (newType) {
    case 'text':
      return { ...base, type: 'text' }

    case 'textarea':
      return { ...base, type: 'textarea' }

    case 'number':
      return { ...base, type: 'number' }

    case 'radio':
      return { ...base, type: 'radio', options: ['New Option'] }

    case 'checkbox':
      return { ...base, type: 'checkbox' }
  }
}

function addOptionToRadioQuestion(question: RadioQuestion) {
  question.options.push('New Option')
}
</script>

<template>
  <main class="container py-4">
    <div class="d-flex align-items-center justify-content-between">
      <div class="flex-grow-1 me-3">
        <input type="text" v-model="training.title" class="h1 underline-input" />
        <input type="text" v-model="training.description" class="p underline-input text-muted" />
      </div>
      <div class="text-end flex-shrink-0">
        <p class="mb-1 text-nowrap">
          {{ training.questions.length }}
          question{{ training.questions.length > 1 ? 's' : '' }}
        </p>
      </div>
    </div>

    <hr />

    <form @submit.prevent>
      <div
        v-for="(q, i) in training.questions"
        :key="q.id"
        class="mb-4 p-3 border rounded shadow-sm d-flex justify-content-between"
      >
        <div>
          <label :for="q.id" class="form-label fs-6">
            <input type="text" v-model="q.label" class="underline-input" />
            <span v-if="q.required" class="text-danger">*</span>
          </label>

          <textarea
            v-if="q.type === 'textarea'"
            :id="q.id"
            v-model="answers[q.id] as string"
            class="form-control"
          ></textarea>

          <input
            v-else-if="q.type === 'number'"
            :id="q.id"
            v-model.number="answers[q.id]"
            type="number"
            class="form-control"
            :required="q.required"
          />

          <div v-else-if="q.type == 'radio'">
            <div v-for="(option, i) in q.options" :key="i" class="form-check">
              <input
                type="radio"
                :id="option"
                :value="option"
                :name="q.id"
                v-model="answers[q.id]"
                class="form-check-input pretty-radio"
              />

              <label :for="option" class="form-check-label">
                <input type="text" v-model="q.options[i]" class="underline-input" />
              </label>
            </div>

            <button class="btn btn-bg p-0 m-0 mt-3" @click="addOptionToRadioQuestion(q)">
              <i class="bi bi-plus-lg"></i>
              Add Option
            </button>
          </div>

          <div v-else-if="q.type === 'checkbox'" class="form-check">
            <input :id="q.id" v-model="answers[q.id]" type="checkbox" class="form-check-input" />
          </div>
        </div>

        <!-- Right side -->
        <div>
          <label for="questionType">Question Type</label>
          <select
            class="form-select"
            id="questionType"
            v-model="q.type"
            @change="training.questions[i] = changeQuestionToDifferentType(q, q.type)"
            aria-label="Default select example"
          >
            <option value="radio">Multiple Choice</option>
            <option value="textarea">Text</option>
            <option value="number">Number</option>
          </select>

          <div class="form-check mt-2">
            <input type="checkbox" id="required" class="form-check-input" v-model="q.required" />
            <label for="required" class="form-check-label">Required?</label>
          </div>
        </div>
      </div>

      <div class="text-center text-muted mt-5 mb-5">No questions</div>

      <div class="text-center mt-3">
        <button class="btn btn-primary" @click="addQuestion">
          <i class="bi bi-plus-lg"></i>
          Add Question
        </button>
      </div>
    </form>

    {{ training }}
  </main>
</template>

<style scoped>
.pretty-radio {
  width: 1.3rem;
  height: 1.3rem;
  margin-top: 0.15rem;
  margin-right: 1rem;
  cursor: pointer;

  border: 1px solid var(--bs-secondary);
}

.underline-input {
  border: 0;
  border-bottom: 1px dashed var(--bs-dark);

  padding: 0;
  margin: 0;

  outline: none;
  box-shadow: none;
}

.underline-input:focus {
  border-bottom: 2px solid #000;
  border-bottom-color: var(--bs-primary);
  border-radius: 0;
}
</style>
