import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import { remarkPreview } from './src/plugins/preview.mjs';
import { remarkDocLinks } from './src/plugins/links.mjs';
import { fileURLToPath } from 'node:url';

const DOCS = fileURLToPath(new URL('../docs', import.meta.url));

export default defineConfig({
  markdown: {
    remarkPlugins: [remarkPreview, [remarkDocLinks, { base: DOCS }]],
  },
  integrations: [
    starlight({
      title: 'instrument',
      description:
        'CSS-кит для интерфейсов, которые показывают человеку работу машины',
      defaultLocale: 'root',
      locales: { root: { label: 'Русский', lang: 'ru' } },
      /* Кит подключается той же строкой, что и у потребителя, — из public/,
         без прохода через сборщик. Иначе сайт показывал бы не тот файл,
         который отдаётся из пакета. */
      head: [
        { tag: 'link', attrs: { rel: 'stylesheet', href: '/kit/kit.css' } },
      ],
      customCss: ['./src/styles/docs.css'],
      /* Шапка переопределена ради спрайта иконок: он обязан быть в теле
         каждой страницы, иначе <use href="#i-refresh"> в примерах пуст. */
      components: { Header: './src/components/Header.astro' },
      /* Порядок разделов — по пути читателя, а не по алфавиту: сначала как
         начать, потом основания, потом компоненты, и агентный слой
         отдельным разделом верхнего уровня, потому что он и есть причина
         брать этот кит. */
      sidebar: [
        { label: 'Основания', items: [{ autogenerate: { directory: 'foundations' } }] },
        { label: 'Раскладка', items: [{ autogenerate: { directory: 'layout' } }] },
        { label: 'Действия', items: [{ autogenerate: { directory: 'components/actions' } }] },
        { label: 'Ввод', items: [{ autogenerate: { directory: 'components/inputs' } }] },
        { label: 'Отображение данных', items: [{ autogenerate: { directory: 'components/display' } }] },
        { label: 'Графики', items: [{ autogenerate: { directory: 'components/charts' } }] },
        { label: 'Навигация', items: [{ autogenerate: { directory: 'components/navigation' } }] },
        { label: 'Оверлеи', items: [{ autogenerate: { directory: 'components/overlays' } }] },
        { label: 'Обратная связь', items: [{ autogenerate: { directory: 'components/feedback' } }] },
        { label: 'Агентный слой', items: [{ autogenerate: { directory: 'agent' } }] },
        { label: 'О проекте', items: [{ autogenerate: { directory: 'about' } }] },
      ],
    }),
  ],
});
