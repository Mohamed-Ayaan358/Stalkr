import ButtonModal from "./components/ButtonModal";
import TableContent from "./components/Table";
import axios from "axios";

export default function Home() {
  // Keep the tooltip but maybe dont need the responsive-td class, You would need a truncate class maybe?
  return (
    <main className="flex flex-col items-center justify-between p-4 md:p-24">
      <TableContent />
      <ButtonModal />
    </main>
  );
}
