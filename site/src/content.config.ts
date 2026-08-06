import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';
import { docsSchema } from '@astrojs/starlight/schema';

/* Контент лежит в ../docs, а не в src/content/docs, и это не прихоть:
   страницы обязаны читаться и в репозитории тоже. Относительные ссылки
   вида ./button.md работают на GitHub и на сайте одинаково — а работали
   бы они только в одном месте, если бы контент жил внутри сайта. */
export const collections = {
  docs: defineCollection({
    loader: glob({
      pattern: ['**/*.md', '!internal/**', '!README.md'],
      base: '../docs',
    }),
    /* Своя часть frontmatter: кит описывает группу, зрелость, исходник и
       честную пометку о том, что компонент без скрипта неполон. */
    schema: docsSchema({
      extend: z.object({
        group: z.string().optional(),
        status: z.enum(['stable', 'draft', 'deprecated']).optional(),
        source: z.string().optional(),
        'needs-js': z.string().optional(),
      }),
    }),
  }),
};
