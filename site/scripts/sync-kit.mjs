/* Кит копируется в public/ перед сборкой сайта.

   Копия, а не импорт, по одной причине: сайт обязан показывать РОВНО тот
   CSS, который получит потребитель, — со всеми @import и порядком слоёв,
   а не пропущенный через сборщик. Собранный Vite'ом кит — это уже другой
   файл, и демонстрировал бы он не то, что отдаётся из пакета.

   Каталог назначения в .gitignore: единственный источник — ../src. */

import { readdirSync, copyFileSync, mkdirSync, rmSync } from 'node:fs';
import { join } from 'node:path';

const P = (u) => new URL(u, import.meta.url).pathname.replace(/^\/([A-Z]:)/, '$1');
const SRC = P('../../src/');
const OUT = P('../public/kit/');

rmSync(OUT, { recursive: true, force: true });
mkdirSync(OUT, { recursive: true });

const files = readdirSync(SRC).filter(f => f.endsWith('.css'));
for (const f of files) copyFileSync(join(SRC, f), join(OUT, f));

const ASSETS = P('../../assets/');
mkdirSync(P('../public/assets/'), { recursive: true });
for (const f of readdirSync(ASSETS)) copyFileSync(join(ASSETS, f), join(P('../public/assets/'), f));

/* Спрайт кладётся ещё и внутрь src/, чтобы Astro вшил его на этапе сборки
   через ?raw. Читать его в рантайме компонента нельзя: после бандлинга
   import.meta.url указывает на чанк, а не на исходник, и путь ломается —
   ровно на этом сборка падала. */
mkdirSync(P('../src/generated/'), { recursive: true });
copyFileSync(join(ASSETS, 'sprite.svg'), P('../src/generated/sprite.svg'));

console.log(`кит синхронизирован: ${files.length} файлов CSS`);
