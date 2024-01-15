import ButtonModal from "./components/ButtonModal";
import TableContent from "./components/Table";
import axios from "axios";

export default function Home() {
  return (
    <main className="flex flex-col items-center justify-between p-4 md:p-24">
      <TableContent />
      <ButtonModal />
    </main>
  );
}

// What you shoudl try and do is keep a preview modal
